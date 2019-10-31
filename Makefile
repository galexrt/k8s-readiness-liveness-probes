.PHONY: dockerbuild
dockerbuild:
	docker build -t galexrt/k8s-readiness-liveness-probes:latest .

.PHONY: build
build:
	go build -o application ./cmd/application/
