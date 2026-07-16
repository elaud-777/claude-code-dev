# Plain Kubernetes manifests

Use these if you don't want Helm. For templated/parameterized installs, use the Helm chart in `helm/taskflow` instead.

This setup is self-contained by default: `postgres.yaml` runs Postgres inside the cluster, and the backend/frontend images are built and loaded locally — no external DB or container registry required.

> **DB note:** these manifests run 2 backend replicas, so `DATABASE_URL` must point at Postgres — SQLite is single-writer and not safe to share across pods/replicas. `postgres.yaml` provides an in-cluster Postgres for exactly this reason.

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

For a remote cluster, push to any registry you control instead (e.g. `docker push ghcr.io/<you>/taskflow-backend:latest`) and update the `image:` fields in `backend-deployment.yaml` / `frontend-deployment.yaml` accordingly.

## 2. Edit hosts

Replace `taskflow.example.com` / `api.taskflow.example.com` in `ingress.yaml`, `backend-configmap.yaml` (`CORS_ORIGINS`), and `frontend-configmap.yaml` (`API_BASE`) with your real domains (or `taskflow.local` / `api.taskflow.local` + `/etc/hosts` entries for a local-only cluster).

## 3. Create the namespace, in-cluster Postgres, and secret

```bash
kubectl apply -f namespace.yaml
kubectl apply -f postgres.yaml

kubectl create secret generic taskflow-backend-secret -n taskflow \
  --from-literal=DATABASE_URL='postgres://taskflow:taskflow@taskflow-postgres:5432/taskflow?sslmode=disable' \
  --from-literal=JWT_SECRET="$(openssl rand -hex 32)"
```

(`backend-secret.example.yaml` shows the declarative alternative — do not commit a filled-in copy. Skip `postgres.yaml` and point `DATABASE_URL` at a managed Postgres like Neon instead, if you'd rather not run your own DB.)

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
