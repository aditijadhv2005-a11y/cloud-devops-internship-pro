# Local lab checklist

## Tools
- Go toolchain
- Docker Engine/Desktop
- kubectl
- Kind
- Helm
- Git
- `kubectl-argo-rollouts` plugin

## Create disposable clusters
```bash
kind create cluster --name operator-dev
# or for the progressive delivery lab
kind create cluster --name progressive-dev
kubectl cluster-info
```

Install platform dependencies following the upstream documentation for your selected compatible versions. Istio and Prometheus require more resources than a simple Kind cluster; allocate sufficient CPU/memory. Delete disposable clusters when finished:
```bash
kind delete cluster --name operator-dev
kind delete cluster --name progressive-dev
```

## Troubleshooting commands
```bash
kubectl get pods -A
kubectl get events -A --sort-by=.lastTimestamp
kubectl describe rollout demo -n demo
kubectl get analysisruns -n demo
kubectl logs -n argo-rollouts deploy/argo-rollouts
```
