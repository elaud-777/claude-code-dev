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

{{- define "taskflow.databaseUrl" -}}
{{- if .Values.backend.secret.databaseUrl -}}
{{- .Values.backend.secret.databaseUrl -}}
{{- else if .Values.postgres.enabled -}}
{{- printf "postgres://%s:%s@%s-postgres:5432/%s?sslmode=disable" .Values.postgres.auth.user .Values.postgres.auth.password .Release.Name .Values.postgres.auth.database -}}
{{- else -}}
{{- fail "Set backend.secret.databaseUrl (or backend.secret.existingSecret) when postgres.enabled is false" -}}
{{- end -}}
{{- end -}}
