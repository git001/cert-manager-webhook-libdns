#!/usr/bin/env bash
set -euo pipefail

ORG="${LIBDNS_ORG:-libdns}"
OUTPUT_DIR="${OUTPUT_DIR:-reports/libdns-provider-compat}"
FAIL_ON_INCOMPATIBLE="${FAIL_ON_INCOMPATIBLE:-false}"
DOCS_DIR="${DOCS_DIR:-docs}"
LATEST_DOC="${DOCS_DIR}/01_provider-compat-latest.md"
ARCHIVE_DOC="${DOCS_DIR}/$(date -u +%Y-%m-%d)_provider-compat.md"

mkdir -p "${OUTPUT_DIR}"
mkdir -p "${DOCS_DIR}"

if [[ -f "${LATEST_DOC}" ]]; then
  mv -f "${LATEST_DOC}" "${ARCHIVE_DOC}"
fi

for cmd in gh jq awk sed sort base64; do
  if ! command -v "${cmd}" >/dev/null 2>&1; then
    echo "missing required command: ${cmd}" >&2
    exit 1
  fi
done

if [[ ! -f go.mod ]]; then
  echo "go.mod not found in current directory" >&2
  exit 1
fi

TARGET_VERSION="$(
  awk '
  BEGIN { inreq = 0 }
  {
    sub(/\/\/.*/, "", $0)
    if ($1 == "require" && $2 == "(") { inreq = 1; next }
    if (inreq && $1 == ")") { inreq = 0; next }
    if ($1 == "require" && $2 == "github.com/libdns/libdns") { print $3; exit }
    if (inreq && $1 == "github.com/libdns/libdns") { print $2; exit }
  }
  ' go.mod
)"

if [[ -z "${TARGET_VERSION}" ]]; then
  echo "could not resolve github.com/libdns/libdns version from go.mod" >&2
  exit 1
fi

target_core="${TARGET_VERSION#v}"
target_major="${target_core%%.*}"

repos_file="${OUTPUT_DIR}/repos.txt"
csv_file="${OUTPUT_DIR}/providers.csv"
md_file="${OUTPUT_DIR}/summary.md"
json_file="${OUTPUT_DIR}/providers.json"

: > "${repos_file}"
cat > "${csv_file}" <<'EOF'
repo,module_path,libdns_requirement,status,reason
EOF

fetch_repos() {
  gh api --paginate "/orgs/${ORG}/repos?per_page=100&type=public" --jq '.[].name' > "${repos_file}"

  sort -u "${repos_file}" -o "${repos_file}"
}

extract_required_version() {
  local gomod="$1"
  awk '
  BEGIN { inreq = 0 }
  {
    sub(/\/\/.*/, "", $0)
    if ($1 == "require" && $2 == "(") { inreq = 1; next }
    if (inreq && $1 == ")") { inreq = 0; next }
    if ($1 == "require" && $2 == "github.com/libdns/libdns") { print $3; exit }
    if (inreq && $1 == "github.com/libdns/libdns") { print $2; exit }
  }
  ' "${gomod}"
}

extract_replace_line() {
  local gomod="$1"
  awk '
  BEGIN { inrep = 0 }
  {
    sub(/\/\/.*/, "", $0)
    if ($1 == "replace" && $2 == "(") { inrep = 1; next }
    if (inrep && $1 == ")") { inrep = 0; next }
    if ($1 == "replace" && $2 == "github.com/libdns/libdns") { print $0; exit }
    if (inrep && $1 == "github.com/libdns/libdns") { print $0; exit }
  }
  ' "${gomod}"
}

normalize_semver() {
  local v="${1#v}"
  # Keep only the numeric x.y.z portion for comparison.
  sed -E 's/^([0-9]+(\.[0-9]+){0,2}).*$/\1/' <<<"${v}"
}

semver_le() {
  local a="$1"
  local b="$2"
  [[ "$(printf '%s\n%s\n' "${a}" "${b}" | sort -V | head -n1)" == "${a}" ]]
}

fetch_repos

while IFS= read -r repo; do
  [[ -z "${repo}" ]] && continue
  [[ "${repo}" == "libdns" ]] && continue

  tmp_json="$(mktemp)"
  tmp_err="$(mktemp)"
  if ! gh api "/repos/${ORG}/${repo}/contents/go.mod" > "${tmp_json}" 2> "${tmp_err}"; then
    if grep -Eq 'HTTP 404|Not Found|status code 404' "${tmp_err}"; then
      echo "${repo},,,skipped,no go.mod in repo root" >> "${csv_file}"
    else
      reason="$(tr '\n' ' ' < "${tmp_err}" | sed -E 's/[[:space:]]+/ /g' | sed 's/,/;/g')"
      reason="${reason:0:180}"
      echo "${repo},,,unknown,${reason}" >> "${csv_file}"
    fi
    rm -f "${tmp_json}" "${tmp_err}"
    continue
  fi
  rm -f "${tmp_err}"

  gomod_file="$(mktemp)"
  jq -r '.content' "${tmp_json}" | tr -d '\n' | base64 -d > "${gomod_file}"
  rm -f "${tmp_json}"

  module_path="$(awk '$1 == "module" { print $2; exit }' "${gomod_file}")"
  if [[ "${module_path}" != github.com/libdns/* ]]; then
    echo "${repo},${module_path},,skipped,module path not under github.com/libdns/" >> "${csv_file}"
    rm -f "${gomod_file}"
    continue
  fi

  replace_line="$(extract_replace_line "${gomod_file}")"
  required_version="$(extract_required_version "${gomod_file}")"
  rm -f "${gomod_file}"

  if [[ -n "${replace_line}" ]]; then
    echo "${repo},${module_path},${required_version},needs-manual-check,uses replace for github.com/libdns/libdns" >> "${csv_file}"
    continue
  fi

  if [[ -z "${required_version}" ]]; then
    echo "${repo},${module_path},,needs-manual-check,no direct require on github.com/libdns/libdns" >> "${csv_file}"
    continue
  fi

  required_core="$(normalize_semver "${required_version}")"
  required_major="${required_core%%.*}"

  if [[ -z "${required_major}" || "${required_major}" == "${required_core}" ]]; then
    echo "${repo},${module_path},${required_version},needs-manual-check,could not parse required semver" >> "${csv_file}"
    continue
  fi

  if [[ "${required_major}" != "${target_major}" ]]; then
    echo "${repo},${module_path},${required_version},incompatible,different major than target ${TARGET_VERSION}" >> "${csv_file}"
    continue
  fi

  if semver_le "${required_core}" "${target_core}"; then
    echo "${repo},${module_path},${required_version},compatible,requires <= target ${TARGET_VERSION}" >> "${csv_file}"
  else
    echo "${repo},${module_path},${required_version},incompatible,requires newer than target ${TARGET_VERSION}" >> "${csv_file}"
  fi
done < "${repos_file}"

jq -R -s '
  split("\n")
  | map(select(length > 0))
  | .[1:]
  | map(split(","))
  | map({
      repo: .[0],
      module_path: .[1],
      libdns_requirement: .[2],
      status: .[3],
      reason: .[4]
    })
' "${csv_file}" > "${json_file}"

total="$(awk 'END { print NR - 1 }' "${csv_file}")"
compatible="$(awk -F',' 'NR > 1 && $4 == "compatible" { c++ } END { print c + 0 }' "${csv_file}")"
incompatible="$(awk -F',' 'NR > 1 && $4 == "incompatible" { c++ } END { print c + 0 }' "${csv_file}")"
manual="$(awk -F',' 'NR > 1 && $4 == "needs-manual-check" { c++ } END { print c + 0 }' "${csv_file}")"
skipped="$(awk -F',' 'NR > 1 && $4 == "skipped" { c++ } END { print c + 0 }' "${csv_file}")"
unknown="$(awk -F',' 'NR > 1 && $4 == "unknown" { c++ } END { print c + 0 }' "${csv_file}")"

{
  echo "# libdns provider compatibility report"
  echo
  echo "- Target libdns version in this repo: \`${TARGET_VERSION}\`"
  echo "- Organization scanned: \`${ORG}\`"
  echo "- Total repositories evaluated: \`${total}\`"
  echo "- Compatible: \`${compatible}\`"
  echo "- Incompatible: \`${incompatible}\`"
  echo "- Needs manual check: \`${manual}\`"
  echo "- Skipped: \`${skipped}\`"
  echo "- Unknown/API issues: \`${unknown}\`"
  echo
  echo "## Incompatible repositories"
  echo
  echo "| Repo | Module | Required libdns | Reason |"
  echo "|---|---|---|---|"
  awk -F',' 'NR > 1 && $4 == "incompatible" { printf("| %s | %s | %s | %s |\n", $1, $2, $3, $5) }' "${csv_file}"
  echo
  echo "## Needs manual check"
  echo
  echo "| Repo | Module | Required libdns | Reason |"
  echo "|---|---|---|---|"
  awk -F',' 'NR > 1 && $4 == "needs-manual-check" { printf("| %s | %s | %s | %s |\n", $1, $2, $3, $5) }' "${csv_file}"
} > "${md_file}"

cp "${md_file}" "${LATEST_DOC}"

echo "Wrote:"
echo "  ${csv_file}"
echo "  ${json_file}"
echo "  ${md_file}"
echo "  ${LATEST_DOC}"

if [[ "${FAIL_ON_INCOMPATIBLE}" == "true" && "${incompatible}" != "0" ]]; then
  echo "incompatible providers found: ${incompatible}" >&2
  exit 1
fi
