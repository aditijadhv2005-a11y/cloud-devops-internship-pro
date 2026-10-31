# Progressive Delivery GitOps

GitOps configuration for Argo CD + Argo Rollouts + Istio + Prometheus. The example uses a canary rollout with pauses and metric analysis.

## Before syncing
1. Replace `REPLACE-ME` repository URLs and image references.
2. Set a real hostname in `apps/demo/istio/gateway.yaml` and configure DNS/TLS for your environment.
3. Install Istio, Argo CD, Argo Rollouts, and Prometheus in the cluster using the upstream instructions.
4. Ensure the Prometheus service DNS in `AnalysisTemplate` matches your installation.
5. Confirm the Istio ingress gateway and telemetry are working.

## Local cluster high-level setup
Use a disposable Kind cluster first. Install each component from its official docs and compatible pinned release, then install the Argo Rollouts kubectl plugin. Do not copy/paste unreviewed installation commands into a production cluster.

## Release workflow
- Update `apps/demo/rollout.yaml` to a new immutable image tag (prefer digest pinning in production).
- Commit and push. Argo CD syncs manifests.
- Watch: `kubectl argo rollouts get rollout demo -n demo --watch`
- Observe: `kubectl argo rollouts dashboard` and query Prometheus.
- If AnalysisTemplate fails, the rollout aborts and stable traffic should continue through Istio routing.

## Traffic splitting
Istio `VirtualService` routes a configured percentage to the canary subset. The rollout's `setWeight` steps and `stableService`/`canaryService` names are intended to coordinate with Argo Rollouts' Istio traffic router. Verify this integration against your installed Argo Rollouts version before a live deployment.
