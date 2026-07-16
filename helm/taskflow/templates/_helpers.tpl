{{- define "taskflow.namespace" -}}
{{- .Values.namespace.name | default .Release.Namespace -}}
{{- end -}}

{{- define "taskflow.backendSecretName" -}}
{{- if .Values.backend.secret.existingSecret -}}
{{- .Values.backend.secret.existingSecret -}}
{{- else -}}
{{- printf "%s-backend-secret" .Release.Name -}}
{{- end -}}
{{- end -}}

{{- define "taskflow.labels" -}}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
