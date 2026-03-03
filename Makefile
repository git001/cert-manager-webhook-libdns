.PHONY: build test test-unit clean rendered-manifest container-build container-push check-container-runtime

IMAGE_NAME ?= cert-manager-webhook-libdns
IMAGE_TAG ?= latest
REGISTRY ?= ghcr.io/your-org
GO_CACHE_DIR ?= $(CURDIR)/.gocache
CONTAINER_RUNTIME ?= $(shell if command -v podman >/dev/null 2>&1; then echo podman; elif command -v docker >/dev/null 2>&1; then echo docker; fi)
CONTAINERFILE ?= Containerfile

# Build the webhook binary
build:
	CGO_ENABLED=0 GOCACHE=$(GO_CACHE_DIR) go build -o webhook -ldflags '-s -w' .

# Run default test suite (unit tests only, no external control-plane dependencies)
test:
	GOCACHE=$(GO_CACHE_DIR) go test -v ./...

# Run unit tests only (same behavior as `make test`)
test-unit:
	GOCACHE=$(GO_CACHE_DIR) go test -v -short ./...

# Clean build artifacts
clean:
	rm -f webhook
	rm -rf _out/

# Generate rendered Helm manifest
rendered-manifest:
	helm template libdns-webhook ./deploy/libdns-webhook \
		--namespace cert-manager \
		> _out/rendered-manifest.yaml

# Ensure a container runtime is available (prefer podman, fallback docker)
check-container-runtime:
	@if [ -z "$(CONTAINER_RUNTIME)" ]; then \
		echo "No container runtime found. Install podman or docker."; \
		exit 1; \
	fi

# Build container image
container-build: check-container-runtime
	$(CONTAINER_RUNTIME) build --file $(CONTAINERFILE) -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Build and push container image
container-push: container-build
	$(CONTAINER_RUNTIME) tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	$(CONTAINER_RUNTIME) push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Install dependencies
deps:
	GOCACHE=$(GO_CACHE_DIR) go mod download
	GOCACHE=$(GO_CACHE_DIR) go mod tidy

# Lint the code
lint:
	GOCACHE=$(GO_CACHE_DIR) golangci-lint run ./...

# Format the code
fmt:
	GOCACHE=$(GO_CACHE_DIR) go fmt ./...

# Verify dependencies
verify:
	GOCACHE=$(GO_CACHE_DIR) go mod verify
