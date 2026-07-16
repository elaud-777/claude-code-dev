# TaskFlow Helm Chart

Templated equivalent of `k8s/*.yaml`. Use this if you manage multiple environments (dev/staging/prod) or want `helm upgrade`/rollback.

This chart is self-contained by default: `postgres.enabled: true` deploys Postgres inside the cluster (an `initContainer` on the backend waits for it via `pg_isready` before starting), and the backend/frontend images default to locally-built, locally-loaded images — no external DB or container registry required.

> **DB note:** `backend.replicaCount` defaults to 2, so the DB must be Postgres, not SQLite — SQLite is single-writer and unsafe to share across replicas. That's exactly what the bundled `postgres.enabled` Postgres is for.

## 1. Build images locally (no registry push needed for a local cluster)

```bash
docker build -t taskflow-backend:local ./backend
docker build -t taskflow-frontend:local ./frontend

# kind:
kind load docker-image taskflow-backend:local taskflow-frontend:local
# minikube:
minikube image load taskflow-backend:local
minikube image load taskflow-frontend:local
```

For a remote cluster, push to any registry you control (e.g. `ghcr.io/<you>/taskflow-backend`) and override `backend.image.repository`/`.tag` (and the frontend equivalents) below.

## 2. Install (fully local — in-cluster Postgres, local images)

```bash
helm install taskflow ./helm/taskflow \
  --set backend.corsOrigins=https://taskflow.yourdomain.com \
  --set frontend.apiBase=https://api.taskflow.yourdomain.com \
  --set ingress.frontendHost=taskflow.yourdomain.com \
  --set ingress.apiHost=api.taskflow.yourdomain.com
```

`backend.secret.databaseUrl` is left blank by default, which auto-derives a connection string to the bundled in-cluster Postgres (`<release>-postgres:5432`) from `postgres.auth.*`. Override `postgres.auth.password` for anything beyond a quick local test.

Or copy `values.yaml` to `my-values.yaml`, edit it, and run `helm install taskflow ./helm/taskflow -f my-values.yaml`.

## 3. Using a managed Postgres (Neon, etc.) instead — optional

```bash
kubectl create secret generic taskflow-backend-secret -n taskflow \
  --from-literal=DATABASE_URL='postgres://user:pass@host/db?sslmode=require' \
  --from-literal=JWT_SECRET="$(openssl rand -hex 32)"

helm install taskflow ./helm/taskflow \
  --set postgres.enabled=false \
  --set backend.secret.existingSecret=taskflow-backend-secret \
  ...
```

Setting `existingSecret` keeps real credentials out of `helm get values`/`helm history` regardless of whether you use the bundled Postgres or an external one.

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
| `postgres.enabled` | `true` | Deploys Postgres in-cluster; set `false` to use `existingSecret`/`databaseUrl` pointing elsewhere |
| `backend.replicaCount` | `2` | Requires Postgres (in-cluster or external), not SQLite, if > 1 |
| `backend.image.repository` / `.tag` | `taskflow-backend` / `local` | Assumes a locally-built, locally-loaded image; override for a remote registry |
| `backend.corsOrigins` | `https://taskflow.example.com` | Comma-separated allowed origins |
| `backend.secret.databaseUrl` | `""` (auto-derived from `postgres.*`) | Set explicitly to point at a different Postgres |
| `backend.secret.existingSecret` | `""` | Name of a pre-created Secret with `DATABASE_URL`/`JWT_SECRET`; overrides `backend.secret.databaseUrl`/`jwtSecret` and the auto-derived in-cluster URL |
| `frontend.apiBase` | `https://api.taskflow.example.com` | Injected as `window.API_BASE` via ConfigMap |
| `ingress.enabled` | `true` | Requires an ingress controller (e.g. ingress-nginx) |
| `ingress.frontendHost` / `ingress.apiHost` | `taskflow.example.com` / `api.taskflow.example.com` | Separate hosts so the frontend and backend don't need path rewriting |
| `ingress.tls.enabled` | `false` | Set `true` + `ingress.tls.secretName` for HTTPS (e.g. via cert-manager) |

## What was validated in this environment

- ✅ `helm lint ./helm/taskflow` — passes
- ✅ `helm template taskflow ./helm/taskflow` — renders correctly: default (postgres enabled, 3 Deployments/Services + 1 PVC + `wait-for-postgres` initContainer), `postgres.enabled=false` + `existingSecret` (postgres resources correctly omitted), and the `fail()` guard when neither `databaseUrl` nor `existingSecret` nor `postgres.enabled` is set
- ⚠️ No Kubernetes cluster was available here, so `helm install`/`upgrade` against a real API server, and the plain manifests in `../k8s/`, still need a final check in your actual cluster.
