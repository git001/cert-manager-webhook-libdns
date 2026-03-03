# Build dependencies stage
FROM golang:1.26-alpine AS build_deps

RUN apk add --no-cache git

WORKDIR /workspace

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Build stage
FROM build_deps AS build

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -o webhook -ldflags '-s -w -extldflags "-static"' .

# Runtime stage
FROM alpine:3.21

ARG OCI_SOURCE="https://github.com/git001/cert-manager-webhook-libdns"
ARG OCI_REVISION=""
ARG OCI_VERSION=""

LABEL org.opencontainers.image.source="${OCI_SOURCE}" \
      org.opencontainers.image.revision="${OCI_REVISION}" \
      org.opencontainers.image.version="${OCI_VERSION}"

RUN apk add --no-cache ca-certificates bash curl

COPY --from=build /workspace/webhook /usr/local/bin/webhook

# Run as non-root user
USER 1001

ENTRYPOINT ["webhook"]
