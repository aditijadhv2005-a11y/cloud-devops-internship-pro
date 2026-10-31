# Security and operational notes

- Never commit kubeconfigs, cloud keys, registry tokens, GitHub tokens, or plaintext production secrets.
- Use GitHub Actions OIDC for cloud publishing/deployment where supported; prefer short-lived credentials.
- Pin action versions to reviewed full commit SHAs before production use. Tags in examples are readable placeholders, not a supply-chain guarantee.
- Use non-root containers, read-only root filesystems where possible, dropped Linux capabilities, seccomp, resource requests/limits, and restricted Pod Security.
- Restrict Kubernetes RBAC to required resources. The operator needs permissions only for its CRD, Deployments, Services, Events, and status updates.
- Scope Istio ingress exposure and TLS to a real, controlled domain. The sample gateway host is a placeholder.
- Treat Prometheus metrics and logs as potentially sensitive. Set retention and access controls.
- Canary analysis thresholds must be calibrated against traffic volume. A low-traffic service may not generate statistically useful metrics.
- Test rollback and recovery in a disposable environment before relying on them.
- Review upstream release notes and compatibility matrices before pinning versions for a production cluster.
