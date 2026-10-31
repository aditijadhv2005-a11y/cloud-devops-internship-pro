# Cloud & DevOps Internship — Professional Project Bundle

This bundle contains three separate, repository-ready projects for the two advanced internship deliverables:

1. **`app-lifecycle-operator/`** — Go + Kubebuilder-style Kubernetes operator managing a `ManagedApp` custom resource.
2. **`progressive-delivery-app/`** — small HTTP application and container build context used as a release workload.
3. **`progressive-delivery-gitops/`** — Argo CD + Argo Rollouts + Istio + Prometheus GitOps configuration for progressive delivery.

## Architecture and ownership boundaries

- The operator owns only resources created from `ManagedApp` (the sample operator reconciles a Deployment and Service).
- Argo Rollouts owns the canary workload in the progressive-delivery example. **Do not let the operator and Argo Rollouts manage the same Deployment.**
- The GitOps repository owns cluster configuration and release manifests; the application repository owns application source and image build.
- The sample uses local Kind as the first target. Production deployment requires cluster-specific domains, TLS, registry credentials, policies, and observability sizing.

## Repository layout

```text
cloud-devops-internship-pro/
├── app-lifecycle-operator/
├── progressive-delivery-app/
├── progressive-delivery-gitops/
└── docs/
```

## Recommended execution order

1. Read `docs/IMPLEMENTATION_PLAN.md` and `docs/SECURITY.md`.
2. Install Go, Docker, kubectl, Kind, Helm, and `controller-gen`/Kubebuilder tools.
3. Run the operator unit tests and build locally.
4. Create a Kind cluster and install the operator.
5. Apply the `ManagedApp` example and inspect status/events.
6. Build and publish the demo application image to your own registry.
7. Install Argo CD, Argo Rollouts, Istio, and Prometheus using the upstream installation instructions for compatible versions.
8. Update image references, hostnames, and GitOps repo URLs; sync the GitOps app and observe the canary.

## Important honesty note

These are structured starter repositories with real implementation files, manifests, tests, workflows, and runbooks. They have **not been run against your machine or a live Kubernetes cluster from this chat**. Pin/check compatible upstream versions before a real cluster deployment. Treat the CI and Kubernetes manifests as a baseline to validate in your environment, not as proof that a production rollout has already succeeded.
