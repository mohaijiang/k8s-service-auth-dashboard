{{/*
Expand the name of the chart.
*/}}
{{- define "k8s-service-auth-dashboard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "k8s-service-auth-dashboard.fullname" -}}
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
{{- define "k8s-service-auth-dashboard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "k8s-service-auth-dashboard.labels" -}}
helm.sh/chart: {{ include "k8s-service-auth-dashboard.chart" . }}
{{ include "k8s-service-auth-dashboard.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "k8s-service-auth-dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k8s-service-auth-dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "k8s-service-auth-dashboard.backend.labels" -}}
{{ include "k8s-service-auth-dashboard.labels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "k8s-service-auth-dashboard.frontend.labels" -}}
{{ include "k8s-service-auth-dashboard.labels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "k8s-service-auth-dashboard.backend.selectorLabels" -}}
{{ include "k8s-service-auth-dashboard.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "k8s-service-auth-dashboard.frontend.selectorLabels" -}}
{{ include "k8s-service-auth-dashboard.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Backend service account name
*/}}
{{- define "k8s-service-auth-dashboard.backend.serviceAccountName" -}}
{{- if .Values.backend.serviceAccount.create }}
{{- default (printf "%s-backend" (include "k8s-service-auth-dashboard.fullname" .)) .Values.backend.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.backend.serviceAccount.name }}
{{- end }}
{{- end }}
