{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cluster-registry.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cluster-registry.fullname" -}}
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
{{- define "cluster-registry.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}


{{- define "cluster-registry-controller.fullname" -}}
{{ include "cluster-registry.fullname" . }}-controller
{{- end }}

{{- define "cluster-registry-controller.name" -}}
{{ include "cluster-registry.name" . }}-controller
{{- end }}

{{- define "cluster-registry-controller.labels" }}
app: {{ include "cluster-registry-controller.fullname" . }}
app.kubernetes.io/name: {{ include "cluster-registry-controller.name" . }}
helm.sh/chart: {{ include "cluster-registry.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | replace "+" "_" }}
app.kubernetes.io/component: cluster-registry-controller
app.kubernetes.io/part-of: {{ include "cluster-registry.name" . }}
{{- end }}

{{- define "cluster-registry-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cluster-registry-controller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
