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
When apps are created in the org namespace add a cluster prefix.
*/}}
{{- define "app.name" -}}
{{- if ne .cluster .ns -}}
{{- printf "%s-%s" .cluster .app -}}
{{- else -}}
{{- .app -}}
{{- end -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "labels.common" -}}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
giantswarm.io/managed-by: {{ .Release.Name | quote }}
giantswarm.io/cluster: {{ .Values.clusterID | quote }}
giantswarm.io/organization: {{ .Values.organization | quote }}
application.giantswarm.io/team: {{ index .Chart.Annotations "application.giantswarm.io/team" | quote }}
helm.sh/chart: {{ include "chart" . | quote }}
{{- end -}}

{{/*
Helpers for App CR re-creation logic. This is a pre-upgrade hook that deletes an existing App CR before a new one is applied.
*/}}
{{- define "re-create-app-cr-hook.updateLogic" -}}
{{- printf "%s-%s" ( default .Chart.Name .Values.nameOverride | trunc 63 ) "update-logic" | replace "+" "_" | trimSuffix "-" -}}
{{- end -}}

{{/*
The team label here only applies to CRD install resources and is to keep the linter happy until we can set this on upstream charts.
*/}}
{{- define "re-create-app-cr-hook.annotations" -}}
"helm.sh/hook": "pre-upgrade"
"helm.sh/hook-delete-policy": "before-hook-creation,hook-succeeded"
application.giantswarm.io/team: {{ index .Chart.Annotations "application.giantswarm.io/team" | quote }}
{{- end -}}

{{- define "re-create-app-cr-hook.shortName" -}}
{{- printf "%s" "re-create-app-cr" -}}
{{- end -}}

{{/* Name template for use in hook resources */}}
{{- define "re-create-app-cr-hook.uniqueName" -}}
{{- printf "%s-%s" .Release.Name ( include "re-create-app-cr-hook.shortName" . ) | replace "+" "_" | trimSuffix "-" -}}
{{- end -}}

{{/* Selector labels for hook resources */}}
{{- define "re-create-app-cr-hook.selectorLabels" -}}
app.kubernetes.io/name: "{{ include "re-create-app-cr-hook.uniqueName" . }}"
app.kubernetes.io/instance: "{{ include "re-create-app-cr-hook.uniqueName" . }}"
{{- end -}}

{{/* Combined labels for hook resources */}}
{{- define "re-create-app-cr-hook.allLabels" -}}
{{ include "labels.common" . }}
app.kubernetes.io/component: "{{ include "re-create-app-cr-hook.shortName" . }}"
{{- end -}}
