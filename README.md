# cert-manager-webhook-libdns

A cert-manager webhook solver for DNS-01 challenges using [libdns](https://github.com/libdns/libdns) providers. This allows you to use any DNS provider that has a libdns implementation for ACME DNS-01 certificate validation.

## Supported DNS Providers

Currently enabled providers:

| Provider | Version | Credential Keys | Documentation |
|----------|---------|-----------------|---------------|
| **Alidns** | v1.0.6-beta.3 | `access_key_id`, `access_key_secret` | [libdns/alidns](https://github.com/libdns/alidns) |
| **Cloudflare** | latest | `api_token` | [libdns/cloudflare](https://github.com/libdns/cloudflare) |
| **deSEC** | v1.0.1 | `api_token` | [libdns/desec](https://github.com/libdns/desec) |
| **Hetzner** | v2.0.1 | `api_token` | [libdns/hetzner](https://github.com/libdns/hetzner) |
| **Linode** | v0.5.0 | `api_token` | [libdns/linode](https://github.com/libdns/linode) |
| **OVH** | v1.1.0 | `endpoint`, `application_key`, `application_secret`, `consumer_key` | [libdns/ovh](https://github.com/libdns/ovh) |
| **Route53** | v1.6.0 | `access_key_id`, `secret_access_key`, `region` | [libdns/route53](https://github.com/libdns/route53) |

Additional providers can be added easily - see [Adding a New Provider](#adding-a-new-provider) section.

### Other Compatible Providers (libdns v1.1.x)

These providers support the latest libdns API and can be added:

| Provider | Version | Module Path | Credential Keys |
|----------|---------|-------------|-----------------|
| **GoDaddy** | latest | `github.com/libdns/godaddy` | `api_token`, `api_secret` |

## Prerequisites

- Kubernetes cluster (tested on OpenShift/CRC)
- [cert-manager](https://cert-manager.io/) v1.0+ installed
- [Helm](https://helm.sh/) v3+ for deployment

## Installation

### 1. Install cert-manager

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.19.3/cert-manager.yaml
```

For DNS-01 challenges, configure cert-manager to use external DNS servers:

```bash
kubectl patch deployment cert-manager -n cert-manager --type='json' \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--dns01-recursive-nameservers=8.8.8.8:53,1.1.1.1:53"},
       {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--dns01-recursive-nameservers-only"}]'
```

### 2. Build and Push the Webhook Image

```bash
# Build the image
podman build -t your-registry/libdns-webhook:v1 .

# Push to your registry
podman push your-registry/libdns-webhook:v1
```

### 3. Deploy the Webhook

```bash
helm install libdns-webhook ./deploy/libdns-webhook \
  --namespace cert-manager \
  --set image.repository=your-registry/libdns-webhook \
  --set image.tag=v1 \
  --set groupName=acme.yourdomain.com
```

**Important:** The `groupName` must be a unique domain you control to avoid conflicts with other webhooks.

### Command Line Options

The webhook binary supports the following options:

```bash
# List compiled-in DNS providers
./webhook --list-providers

# Example output:
# Compiled-in DNS providers:
#   - alidns
#   - cloudflare
#   - desec
#   - hetzner
#   - linode
#   - ovh
#   - route53
```

This is useful to verify which providers are available in your build.

## Configuration

### Create DNS Provider Credentials Secret

Create a Kubernetes Secret containing your DNS provider credentials:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dns-provider-credentials
  namespace: cert-manager
type: Opaque
stringData:
  api_token: "your-api-token-here"
```

For Route53:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: route53-credentials
  namespace: cert-manager
type: Opaque
stringData:
  access_key_id: "AKIAIOSFODNN7EXAMPLE"
  secret_access_key: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  region: "us-east-1"
```

### Create a ClusterIssuer

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod-account
    solvers:
      - dns01:
          webhook:
            groupName: acme.yourdomain.com  # Must match Helm groupName
            solverName: libdns
            config:
              provider: desec  # or cloudflare, hetzner, route53
              ttl: 3600  # Optional: TTL in seconds (default: 300, deSEC requires 3600)
              secretRef:
                name: dns-provider-credentials
                namespace: cert-manager
              # zone: example.com  # Optional: override auto-detected zone
```

### Request a Certificate

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-certificate
  namespace: default
spec:
  secretName: my-tls-secret
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
    - example.com
    - "*.example.com"
```

## Provider-Specific Notes

### deSEC

- Minimum TTL is 3600 seconds (enforced by deSEC) - set `ttl: 3600` in config
- API to DNS propagation can take up to 2 minutes
- API token can be created at https://desec.io/tokens

### Cloudflare

- Use an API token with `Zone:DNS:Edit` permissions
- Scoped tokens are recommended over global API keys

### Route53

- IAM user needs `route53:ChangeResourceRecordSets` and `route53:ListHostedZones` permissions
- The `region` field is optional (defaults to `us-east-1`)

### Alidns (Alibaba Cloud)

- Create an AccessKey at https://ram.console.aliyun.com/manage/ak
- The `region_id` field is optional (defaults to `cn-hangzhou`)
- For temporary credentials, provide `security_token`

### OVH

- Create API credentials at https://api.ovh.com/createToken/
- Required permissions: `GET /domain/zone/*`, `POST /domain/zone/*`, `PUT /domain/zone/*`, `DELETE /domain/zone/*`
- Endpoint values: `ovh-eu`, `ovh-ca`, `ovh-us`, `kimsufi-eu`, `kimsufi-ca`, `soyoustart-eu`, `soyoustart-ca`

### Linode

- Create an API token at https://cloud.linode.com/profile/tokens
- Token requires `Domains` read/write permissions

### Hetzner

- Create an API token at https://dns.hetzner.com/settings/api-token
- Token requires read/write permissions

## Wildcard + Base Domain Certificates

When requesting both a wildcard (`*.example.com`) and base domain (`example.com`) certificate, the webhook handles creating multiple TXT records at `_acme-challenge.example.com` with different values.

```yaml
dnsNames:
  - example.com      # Needs TXT record with key1
  - "*.example.com"  # Needs TXT record with key2 (same DNS name!)
```

The webhook uses `GetRecords` + `SetRecords` to merge multiple TXT values at the same DNS name.

## Troubleshooting

### Check Webhook Logs

```bash
kubectl logs -n cert-manager -l app.kubernetes.io/name=libdns-webhook
```

### Check Challenge Status

```bash
kubectl get challenges -A
kubectl describe challenge <challenge-name> -n <namespace>
```

### Common Issues

**1. "DNS record not yet propagated"**
- Wait for DNS TTL to expire (can be up to 1 hour for providers with high minimum TTL)
- Verify the TXT record exists: `dig TXT _acme-challenge.yourdomain.com @8.8.8.8`

**2. "failed to get secret"**
- Ensure the credentials secret exists in the correct namespace
- Check RBAC permissions for the webhook service account

**3. Pod fails to start on OpenShift**
- The Helm chart is configured for OpenShift SCC compatibility
- Do not set `runAsUser` or `fsGroup` explicitly

**4. "unknown DNS provider"**
- Check that the provider name in the ClusterIssuer config matches a registered provider
- Currently available: `alidns`, `cloudflare`, `desec`, `hetzner`, `linode`, `ovh`, `route53`
- Use `--list-providers` to see compiled-in providers

## Development

### Adding a New Provider

This example shows how to add the **Hetzner** provider (v2.0.1) which supports the latest libdns v1.1.1.

#### Step 1: Check Provider Compatibility

Ensure the libdns provider supports libdns v1.1.1. Check the provider's `go.mod`:

```bash
curl -sL "https://raw.githubusercontent.com/libdns/hetzner/v2.0.1/go.mod" | grep libdns
# Should show: github.com/libdns/libdns v1.1.1
```

**Compatible providers (as of Feb 2026):**
| Provider | Version | Module Path |
|----------|---------|-------------|
| alidns | v1.0.6-beta.3 | `github.com/libdns/alidns` |
| cloudflare | latest | `github.com/libdns/cloudflare` |
| desec | v1.0.1 | `github.com/libdns/desec` |
| godaddy | latest | `github.com/libdns/godaddy` |
| hetzner | v2.0.1 | `github.com/libdns/hetzner/v2` |
| linode | v0.5.0 | `github.com/libdns/linode` |
| ovh | v1.1.0 | `github.com/libdns/ovh` |
| route53 | v1.6.0 | `github.com/libdns/route53` |

> **Note:** Providers with v2+ use a different module path (e.g., `/v2` suffix).

#### Step 2: Create the Provider File

Create `providers/hetzner.go`:

```go
package providers

import (
	"fmt"

	hetzner "github.com/libdns/hetzner/v2"  // Note: v2 module path
)

func init() {
	Register("hetzner", NewHetznerProvider)
}

// NewHetznerProvider creates a Hetzner DNS provider
//
// Required credentials:
//   - api_token: Hetzner DNS API token
func NewHetznerProvider(config ProviderConfig) (DNSProvider, error) {
	apiToken := config.Credentials["api_token"]
	if apiToken == "" {
		return nil, fmt.Errorf("hetzner: api_token is required")
	}

	return &hetzner.Provider{
		AuthAPIToken: apiToken,
	}, nil
}
```

#### Step 3: Add the Dependency

Add the provider to `go.mod`:

```go
require (
    // ... existing dependencies ...
    github.com/libdns/hetzner/v2 v2.0.1
)
```

Or use `go get`:

```bash
go get github.com/libdns/hetzner/v2@v2.0.1
```

#### Step 4: Update Dependencies

Run `go mod tidy` to update `go.sum`:

```bash
go mod tidy
```

#### Step 5: Rebuild and Deploy

```bash
# Build the image
podman build -t your-registry/libdns-webhook:v2 .

# Push and update deployment
podman push your-registry/libdns-webhook:v2
kubectl set image deployment/libdns-webhook -n cert-manager \
  libdns-webhook=your-registry/libdns-webhook:v2
```

#### Step 6: Configure the Issuer

Create credentials secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: hetzner-credentials
  namespace: cert-manager
type: Opaque
stringData:
  api_token: "your-hetzner-dns-api-token"
```

Reference in ClusterIssuer:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod-account
    solvers:
      - dns01:
          webhook:
            groupName: acme.yourdomain.com
            solverName: libdns
            config:
              provider: hetzner
              ttl: 300  # Hetzner allows lower TTLs than deSEC
              secretRef:
                name: hetzner-credentials
                namespace: cert-manager
```

### Provider Interface Requirements

Each libdns provider must implement the `DNSProvider` interface:

```go
type DNSProvider interface {
    libdns.RecordAppender  // AppendRecords(ctx, zone, records)
    libdns.RecordDeleter   // DeleteRecords(ctx, zone, records)
    libdns.RecordGetter    // GetRecords(ctx, zone)
    libdns.RecordSetter    // SetRecords(ctx, zone, records)
}
```

The webhook uses:
- `GetRecords` + `SetRecords` for **Present()** (to merge multiple TXT values)
- `GetRecords` + `SetRecords` or `DeleteRecords` for **CleanUp()**

### Running Tests

```bash
go test ./...
```

## Architecture

```
                                    ┌─────────────────┐
                                    │   cert-manager  │
                                    └────────┬────────┘
                                             │
                                             │ DNS-01 Challenge
                                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    libdns-webhook                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   Present()  │  │   CleanUp()  │  │ Initialize() │          │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘          │
│         │                 │                                     │
│         ▼                 ▼                                     │
│  ┌─────────────────────────────────┐                           │
│  │      Provider Registry          │                           │
│  │  ┌──────────┐ ┌──────────┐     │                           │
│  │  │cloudflare│ │  desec   │ ... │                           │
│  │  └──────────┘ └──────────┘     │                           │
│  └─────────────────────────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
                                             │
                                             │ libdns interface
                                             ▼
                                    ┌─────────────────┐
                                    │   DNS Provider  │
                                    │      API        │
                                    └─────────────────┘
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
