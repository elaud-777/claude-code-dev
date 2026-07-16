# Plain Kubernetes manifests

Use these if you don't want Helm. For templated/parameterized installs, use the Helm chart in `helm/taskflow` instead.

> **DB note:** these manifests run 2 backend replicas, so `DATABASE_URL` must point at Postgres (Neon) — SQLite (the local/Docker Compose default) is single-writer and not safe to share across pods/replicas.

## 1. Build & push images

```bash
docker build -t ghcr.io/<you>/taskflow-backend:latest ./backend
docker build -t ghcr.io/<you>/taskflow-frontend:latest ./frontend
docker push ghcr.io/<you>/taskflow-backend:latest
docker push ghcr.io/<you>/taskflow-frontend:latest
```

Update the `image:` fields in `backend-deployment.yaml` / `frontend-deployment.yaml` (or `helm/taskflow/values.yaml`) to match.

## 2. Edit hosts

Replace `taskflow.example.com` / `api.taskflow.example.com` in `ingress.yaml`, `backend-configmap.yaml` (`CORS_ORIGINS`), and `frontend-configmap.yaml` (`API_BASE`) with your real domains.

## 3. Create the secret

```bash
kubectl apply -f namespace.yaml
kubectl create secret generic taskflow-backend-secret -n taskflow \
  --from-literal=DATABASE_URL='postgres://user:pass@host/db?sslmode=require' \
  --from-literal=JWT_SECRET="$(openssl rand -hex 32)"
```

(`backend-secret.example.yaml` shows the declarative alternative — do not commit a filled-in copy.)

## 4. Apply everything else

```bash
kubectl apply -f backend-configmap.yaml
kubectl apply -f backend-deployment.yaml
kubectl apply -f backend-service.yaml
kubectl apply -f frontend-configmap.yaml
kubectl apply -f frontend-deployment.yaml
kubectl apply -f frontend-service.yaml
kubectl apply -f ingress.yaml
```

## 5. Verify

```bash
kubectl -n taskflow get pods
kubectl -n taskflow logs deploy/taskflow-backend
curl https://api.taskflow.example.com/health
```

Requires an ingress controller (e.g. ingress-nginx) already installed in the cluster.
