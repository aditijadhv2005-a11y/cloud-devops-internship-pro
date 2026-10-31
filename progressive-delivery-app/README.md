# Progressive Delivery Demo App

A tiny Go HTTP service with `/`, `/healthz`, and `/version` endpoints. It returns a configurable version and supports a test-only failure mode for demonstrating rollback.

```bash
go test ./...
go run .
curl http://localhost:8080/healthz
```

Build a container:
```bash
docker build -t YOUR_REGISTRY/progressive-demo:0.1.0 .
docker push YOUR_REGISTRY/progressive-demo:0.1.0
```
Set `APP_VERSION` to the release identifier. Set `FAIL_HEALTH=true` only in a controlled rollback exercise; do not enable it in normal deployments.
