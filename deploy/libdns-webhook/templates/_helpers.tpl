{{/*
Expand the name of the chart.
*/}}
{{- define "libdns-webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "libdns-webhook.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "libdns-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "libdns-webhook.labels" -}}
helm.sh/chart: {{ include "libdns-webhook.chart" . }}
{{ include "libdns-webhook.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "libdns-webhook.selectorLabels" -}}
app.kubernetes.io/name: {{ include "libdns-webhook.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "libdns-webhook.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "libdns-webhook.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Self-signed issuer name
*/}}
{{- define "libdns-webhook.selfSignedIssuer" -}}
{{ include "libdns-webhook.fullname" . }}-selfsign
{{- end }}

{{/*
CA issuer name
*/}}
{{- define "libdns-webhook.caIssuer" -}}
{{ include "libdns-webhook.fullname" . }}-ca
{{- end }}

{{/*
CA certificate name
*/}}
{{- define "libdns-webhook.caCertificate" -}}
{{ include "libdns-webhook.fullname" . }}-ca
{{- end }}

{{/*
Serving certificate name
*/}}
{{- define "libdns-webhook.servingCertificate" -}}
{{ include "libdns-webhook.fullname" . }}-tls
{{- end }}
