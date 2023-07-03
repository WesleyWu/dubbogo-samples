{{/*
Expand the name of the chart.
*/}}
{{- define "app.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
Can we fix line 'if .Values.version.labels.dubbogoAppVersion' if user doesn't want to set app version?
*/}}
{{- define "app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- if .Values.version.labels.dubbogoAppVersion }}
{{- $version := .Values.version.labels.dubbogoAppVersion }}
{{- printf "%s-%s" .Chart.Name $version  | trimSuffix "-" }}
{{- else }}
{{- printf "%s" .Chart.Name }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "app.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

*/}}
{{- define "app.runLabels" -}}
run: {{ include "app.name" . }}
{{- end }}
