# Canary rollout runbook

## Pre-flight
- Confirm `kubectl config current-context` points to a disposable/test cluster.
- Confirm Argo CD application is healthy and synced.
- Confirm Istio ingress, Prometheus, and Argo Rollouts controllers are healthy.
- Confirm the new image exists and is immutable; preferably use an image digest.
- Check baseline request success rate and p95 latency.

## Start a release
1. Update the image tag/digest in `apps/demo/rollout.yaml`.
2. Open a pull request and pass manifest validation.
3. Merge after review; wait for Argo CD sync.
4. Watch `kubectl argo rollouts get rollout demo -n demo --watch`.
5. Inspect AnalysisRuns: `kubectl get analysisrun -n demo`.

## Abort / rollback
- If automated analysis fails, verify the rollout is aborted and stable ReplicaSet remains healthy.
- For a manual abort, use `kubectl argo rollouts abort demo -n demo`.
- Revert the Git commit to the last known-good image; let Argo CD reconcile desired state.
- Confirm stable traffic, error rate, and latency recover; record the incident and root cause.

## Evidence to capture for assessment
- Git commit/PR that changes the immutable image.
- Argo CD application synced/healthy view.
- Rollout status at 10%, 25%, 50%, and 100%.
- AnalysisRun success and a controlled failure/abort demonstration.
- Prometheus graph showing success rate and p95 latency.
- Rollback timeline and explanation of why user traffic remained available.

## Caveats
Prometheus metric names depend on your telemetry pipeline. The demo application does not itself expose the Istio/Prometheus metrics queried here; configure Istio telemetry or an appropriate metrics adapter and adjust labels/queries to your actual series before expecting AnalysisTemplate success.
