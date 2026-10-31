.PHONY: test build run manifests install uninstall docker-build
IMG ?= ghcr.io/REPLACE-ME/app-lifecycle-operator:0.1.0

deps:
	go mod tidy

test:
	go test ./... -cover

build:
	go build -o bin/manager ./cmd

run:
	go run ./cmd

manifests:
	controller-gen rbac:roleName=manager-role crd paths="./..." output:crd:artifacts:config=config/crd/bases

install: manifests
	kubectl apply -f config/crd/bases

uninstall:
	kubectl delete -f config/crd/bases --ignore-not-found

docker-build:
	docker build -t $(IMG) .
