# Architecture

```text
Developer PR -> CI tests/build -> immutable container image
      |                               |
      +-------- GitOps PR <------------+
                     |
                  Argo CD
                     |
               Argo Rollouts
             /                 \
      stable ReplicaSet     canary ReplicaSet
             \                 /
             Istio traffic routing
                     |
                 end users
                     |
           Prometheus metrics/analysis
```

The operator project is a separate platform API example. It is intentionally not the controller for this Rollout workload, avoiding two reconcilers fighting over the same Deployment.
