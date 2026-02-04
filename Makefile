.PHONY: build test clean rendered-manifest docker-build docker-push

IMAGE_NAME ?= cert-manager-webhook-libdns
IMAGE_TAG ?= latest
REGISTRY ?= ghcr.io/your-org

# Kubernetes version for kubebuilder assets
KUBEBUILDER_VERSION ?= 1.28.0

# Build the webhook binary
build:
	CGO_ENABLED=0 go build -o webhook -ldflags '-s -w' .

# Run conformance tests
# Requires TEST_ZONE_NAME environment variable to be set
test:
	@if [ -z "$(TEST_ZONE_NAME)" ]; then \
		echo "TEST_ZONE_NAME must be set to run conformance tests"; \
		exit 1; \
	fi
	go test -v .

# Run unit tests only (no conformance)
test-unit:
	go test -v -short ./...

# Clean build artifacts
clean:
	rm -f webhook
	rm -rf _out/

# Generate rendered Helm manifest
rendered-manifest:
	helm template libdns-webhook ./deploy/libdns-webhook \
		--namespace cert-manager \
		> _out/rendered-manifest.yaml

# Build Docker image
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

# Build and push Docker image
docker-push: docker-build
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Install dependencies
deps:
	go mod download
	go mod tidy

# Lint the code
lint:
	golangci-lint run ./...

# Format the code
fmt:
	go fmt ./...

# Verify dependencies
verify:
	go mod verify
