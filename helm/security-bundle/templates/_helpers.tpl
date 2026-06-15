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
application.giantswarm.io/team: {{ index .Chart.Annotations "io.giantswarm.application.team" | quote }}
cluster.x-k8s.io/cluster-name: {{ .Values.clusterID | quote }}
helm.sh/chart: {{ include "chart" . | quote }}
{{- end -}}

{{- /*
sortedValuesFrom — reproduce the App platform's config merge order as a single
Flux `valuesFrom` list (Flux merges the list top-to-bottom, later wins).

App platform (github.com/giantswarm/app/pkg/values) merges all configMap layers,
then all secret layers, then the secret blob OVER the configMap blob — so a secret
ALWAYS overrides a configMap regardless of priority. Within a kind, extraConfigs are
ordered by `priority` (1-150, default 25) around the user-config slot (100); bundle
sub-apps have no cluster-values layer. Equal-priority extras keep authored order; an
extra at priority 100 precedes the user-config layer (band is half-open: p <= 100).

Sort key `KIND_PRIORITY_SOURCE_INDEX`:
  KIND   0=configMap 1=secret      (all configMaps before all secrets)
  PRIOR  %03d                       (numeric order)
  SOURCE 0=extraConfig 1=userConfig (extra precedes user layer at equal priority)
  INDEX  %03d authored position     (stable tie-break)
Returns a YAML array of {kind,name} dicts.

Arg: dict with keys:
  extraConfigs        list of {kind,name,priority}
  userConfigMapName   name | ""   (slot 100, configMap)
  userSecretName      name | ""   (slot 100, secret)
  root                $ (for tpl-ing extraConfig names)
*/ -}}
{{- define "sortedValuesFrom" -}}
{{- $keyed := dict -}}
{{- $keys := list -}}
{{- if .userConfigMapName -}}
{{- $k := printf "0_%03d_1_000" 100 -}}
{{- $keyed = set $keyed $k (dict "kind" "ConfigMap" "name" .userConfigMapName) -}}
{{- $keys = append $keys $k -}}
{{- end -}}
{{- if .userSecretName -}}
{{- $k := printf "1_%03d_1_000" 100 -}}
{{- $keyed = set $keyed $k (dict "kind" "Secret" "name" .userSecretName) -}}
{{- $keys = append $keys $k -}}
{{- end -}}
{{- range $i, $ec := (.extraConfigs | default list) -}}
{{- $isSecret := eq (lower ($ec.kind | default "configMap")) "secret" -}}
{{- $k := printf "%s_%03d_0_%03d" (ternary "1" "0" $isSecret) (int ($ec.priority | default 25)) $i -}}
{{- $keyed = set $keyed $k (dict "kind" (ternary "Secret" "ConfigMap" $isSecret) "name" (tpl $ec.name $.root)) -}}
{{- $keys = append $keys $k -}}
{{- end -}}
{{- $out := list -}}
{{- range $k := sortAlpha $keys -}}
{{- $out = append $out (get $keyed $k) -}}
{{- end -}}
{{- $out | toYaml -}}
{{- end -}}
