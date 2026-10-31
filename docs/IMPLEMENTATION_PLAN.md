# Implementation plan

## Milestone 1 — Foundation
- Install Go, Git, Docker, kubectl, Kind, Helm, and Make.
- Verify `go version`, `git --version`, `docker version`, `kubectl version --client`, `kind version`, `helm version`.
- Create one GitHub repository per directory; keep secrets out of Git.

## Milestone 2 — AppLifecycle Operator
- Define `ManagedApp` API with replica count, image, port, and resource requests/limits.
- Generate CRD and RBAC markers from Go API types.
- Reconcile Deployment and Service idempotently.
- Report Ready/Progressing/Degraded conditions and observed generation.
- Add unit tests for desired resource generation and reconcile behavior.
- Build/publish an image and deploy to a disposable Kind cluster.

Acceptance checks:
- Applying the sample CR creates a Deployment and Service.
- Changing `spec.replicas` or `spec.image` converges the Deployment.
- Deleting the CR garbage-collects owned resources.
- Invalid fields are rejected by CRD schema validation.

## Milestone 3 — Progressive delivery
- Build the demo app image and push to a registry you control.
- Install Argo CD, Argo Rollouts, Istio, and Prometheus.
- Configure the Istio gateway/virtual service and destination rules for stable/canary subsets.
- Roll out at 10%, 25%, 50%, then 100% traffic; pause between steps for analysis.
- Use Prometheus-backed AnalysisTemplates for request success rate and latency.
- Demonstrate abort/rollback with a deliberate failing image or a test-only failure flag.

Acceptance checks:
- Git commit changes the desired image tag; Argo CD syncs it.
- Rollout exposes stable and canary ReplicaSets and progresses through configured steps.
- Analysis failure aborts the rollout and stable traffic remains available.
- A dashboard or query shows error rate/latency before, during, and after rollout.

## Milestone 4 — Production hardening
- Add policy enforcement, signed images, SBOM, vulnerability thresholds, secret management, network policies, Pod Security Standards, resource quotas, backup/restore, and SLO alerts.
- Use a managed Kubernetes cluster only after local validation and cost review.
- Add a staged environment promotion process; never auto-promote unverified artifacts to production.
