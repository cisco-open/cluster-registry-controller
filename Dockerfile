# Copyright (c) 2021 Banzai Cloud Zrt. All Rights Reserved.

# Build the manager binary
FROM golang:1.15 as builder

ARG GITHUB_ACCESS_TOKEN

ARG GOPROXY=https://proxy.golang.org
ARG GOPRIVATE=github.com/banzaicloud

ENV GOFLAGS="-mod=readonly"

WORKDIR /workspace/

# Copy the Go Modules manifests
COPY ./go.mod /workspace/
COPY ./go.sum /workspace/
RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/"
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
