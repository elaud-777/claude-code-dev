# TaskFlow Helm Chart

Templated equivalent of `k8s/*.yaml`. Use this if you manage multiple environments (dev/staging/prod) or want `helm upgrade`/rollback.

> **DB note:** `backend.replicaCount` defaults to 2, so `databaseUrl` must point at Postgres (Neon) — SQLite is single-writer and unsafe to share across replicas. Set `backend.replicaCount: 1` only if you intentionally stick with SQLite for a small/non-HA deployment.

## 1. Build & push images

```bash
docker build -t ghcr.io/<you>/taskflow-backend:latest ./backend
docker build -t ghcr.io/<you>/taskflow-frontend:latest ./frontend
docker push ghcr.io/<you>/taskflow-backend:latest
docker push ghcr.io/<you>/taskflow-frontend:latest
```

## 2. Create the DB/JWT secret out-of-band (recommended)

```bash
kubectl create namespace taskflow
kubectl create secret generic taskflow-backend-secret -n taskflow \
  --from-literal=DATABASE_URL='postgres://user:pass@host/db?sslmode=require' \
  --from-literal=JWT_SECRET="$(openssl rand -hex 32)"
```

Then reference it via `--set backend.secret.existingSecret=taskflow-backend-secret` (see below). This keeps secrets out of `helm get values` / `helm history`.

## 3. Install

```bash
helm install taskflow ./helm/taskflow \
  --set backend.image.repository=ghcr.io/<you>/taskflow-backend \
  --set frontend.image.repository=ghcr.io/<you>/taskflow-frontend \
  --set backend.secret.existingSecret=taskflow-backend-secret \
  --set backend.corsOrigins=https://taskflow.yourdomain.com \
  --set frontend.apiBase=https://api.taskflow.yourdomain.com \
  --set ingress.frontendHost=taskflow.yourdomain.com \
  --set ingress.apiHost=api.taskflow.yourdomain.com
```

Or copy `values.yaml` to `my-values.yaml`, edit it, and run `helm install taskflow ./helm/taskflow -f my-values.yaml`.

## 4. Upgrade / rollback

```bash
helm upgrade taskflow ./helm/taskflow -f my-values.yaml
helm rollback taskflow 1
```

## 5. Verify

```bash
helm status taskflow -n taskflow
kubectl -n taskflow get pods
curl https://api.taskflow.yourdomain.com/health
```

## Key values

| Key | Default | Notes |
|---|---|---|
| `backend.replicaCount` | `2` | Requires Postgres, not SQLite, if > 1 |
| `backend.image.repository` / `.tag` | `ghcr.io/elaud-777/taskflow-backend` / `latest` | Update to your own registry |
| `backend.corsOrigins` | `https://taskflow.example.com` | Comma-separated allowed origins |
| `backend.secret.existingSecret` | `""` | Name of a pre-created Secret with `DATABASE_URL`/`JWT_SECRET`; overrides `backend.secret.databaseUrl`/`jwtSecret` |
| `frontend.apiBase` | `https://api.taskflow.example.com` | Injected as `window.API_BASE` via ConfigMap |
| `ingress.enabled` | `true` | Requires an ingress controller (e.g. ingress-nginx) |
| `ingress.frontendHost` / `ingress.apiHost` | `taskflow.example.com` / `api.taskflow.example.com` | Separate hosts so the frontend and backend don't need path rewriting |
| `ingress.tls.enabled` | `false` | Set `true` + `ingress.tls.secretName` for HTTPS (e.g. via cert-manager) |

## What was validated in this environment

- ✅ `helm lint ./helm/taskflow` — passes
- ✅ `helm template taskflow ./helm/taskflow` — renders correctly, including the `existingSecret` conditional (Secret template correctly omitted when set)
- ⚠️ No Kubernetes cluster was available here, so `helm install`/`upgrade` against a real API server, and the plain manifests in `../k8s/`, still need a final check in your actual cluster.
