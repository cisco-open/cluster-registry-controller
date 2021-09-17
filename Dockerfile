# Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

# Build the manager binary
FROM golang:1.17 as builder

ARG GITHUB_ACCESS_TOKEN

ARG GOPROXY="https://proxy.golang.org,direct"
ENV GOPROXY="${GOPROXY}"
ENV GOPRIVATE='github.cisco.com,github.com/banzaicloud'
ENV GONOPROXY='gopkg.in,go.uber.org'
ENV GOFLAGS="-mod=readonly"



ENV GOFLAGS="-mod=readonly"

WORKDIR /workspace/

# Copy the Go Modules manifests
COPY ./go.mod /workspace/
COPY ./go.sum /workspace/
RUN if [ -n "${GITHUB_ACCESS_TOKEN}" ]; then \
      git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"; \
    fi
RUN go mod download

COPY ./ /workspace/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on make binary

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]
