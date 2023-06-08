# Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
ARG GID=1000
ARG UID=1000

# Build the manager binary
FROM golang:1.18.0 as builder
ARG GITHUB_ACCESS_TOKEN
ARG GID
ARG UID

# Create user and group
RUN groupadd -g ${GID} appgroup && \
    useradd -u ${UID} --gid appgroup appuser

ARG GOPROXY="https://proxy.golang.org,direct"
ENV GOPROXY="${GOPROXY}"
ENV GOPRIVATE='github.com/cisco-open,github.com/banzaicloud'
ENV GONOPROXY='gopkg.in,go.uber.org'
ENV GOFLAGS="-mod=readonly"

WORKDIR /workspace/

# Copy the Go Modules manifests
COPY ./go.mod /workspace/
COPY ./go.sum /workspace/
# Copy the API Go Modules manifests
COPY api/go.mod api/go.mod
COPY api/go.sum api/go.sum
RUN if [ -n "${GITHUB_ACCESS_TOKEN}" ]; then \
      git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"; \
    fi
RUN go mod download

COPY ./ /workspace/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on make binary

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ARG GID
ARG UID
WORKDIR /
COPY --from=builder /workspace/bin/manager .
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER ${UID}:${GID}

ENTRYPOINT ["/manager"]
