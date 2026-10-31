# Platform installation references

Install and version these platform dependencies separately from the app manifests:

- Argo CD — GitOps reconciliation
- Argo Rollouts — progressive delivery controller and kubectl plugin
- Istio — service mesh and traffic management
- Prometheus — metric provider for rollout analysis

Use official upstream installation documentation and select versions that are mutually compatible. Avoid installing all platform components into a production cluster from a single unreviewed shell script.
