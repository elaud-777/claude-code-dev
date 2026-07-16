// Runtime config, loaded before js/api.js. Overridden per-deployment:
// - Docker/local: default below (backend on localhost:8000) is used as-is.
// - Kubernetes: a ConfigMap replaces this file with the cluster-specific API host.
window.API_BASE = window.API_BASE || "http://localhost:8000";
