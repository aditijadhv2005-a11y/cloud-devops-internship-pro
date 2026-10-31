# AppLifecycle Operator

A compact Go controller-runtime operator that reconciles `ManagedApp` resources into Kubernetes Deployments and Services.

## Prerequisites
- Go 1.23+
- Docker
- kubectl and Kind for local integration
- controller-gen (for CRD/RBAC generation)

## Quick start
```bash
make deps
make test
make build
make manifests
make install   # requires a configured Kubernetes context
make run       # run controller locally against current context
# in another terminal
kubectl apply -f config/samples/platform_v1alpha1_managedapp.yaml
kubectl get managedapps
kubectl get deploy,svc
```

`make install` installs the CRD. For an isolated cluster:
```bash
kind create cluster --name operator-dev
kubectl cluster-info --context kind-operator-dev
```

## API example
See `config/samples/platform_v1alpha1_managedapp.yaml`. The operator deliberately owns only its generated Deployment and Service. It does not own Argo Rollouts resources.

## Limitations to address before production
- Add integration/envtest coverage for watches, deletion, and status transitions.
- Extend API validation and admission policy to your workload contract.
- Add leader election for multi-replica deployment and configure health/readiness probes.
- Review generated RBAC and restrict it further for your environment.
