# TEMPLATE-APP: This is set as a reasonable default, feel free to change.

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "labels.common" -}}
{{ include "labels.selector" . }}
app-operator.giantswarm.io/version: 5.5.1
app.giantswarm.io/branch: {{ .Chart.Annotations.branch | replace "#" "-" | replace "/" "-" | replace "." "-" | trunc 63 | trimSuffix "-" | quote }}
app.giantswarm.io/commit: {{ .Chart.Annotations.commit | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
helm.sh/chart: {{ include "chart" . | quote }}
giantswarm.io/service-type: {{ .Values.serviceType }}
giantswarm.io/cluster: {{ .Values.clusterName | quote }}
giantswarm.io/organization: {{ .Values.organization | quote }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "labels.selector" -}}
app.kubernetes.io/name: {{ include "name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end -}}

{{- define "kubeconfig" -}}
kubeConfig:
  context:
    name: {{ .Values.clusterName }}-admin@{{ .Values.clusterName }}
  inCluster: false
  secret:
    name: {{ .Values.clusterName }}-kubeconfig
    namespace: {{ .Release.Namespace }}
{{- end -}}

{{- define "config" -}}
config:
  configMap:
    name: {{ .Values.clusterName }}-cluster-values
    namespace: {{ .Release.Namespace }}
{{- end -}}
