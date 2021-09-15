#!/bin/bash

set -x
set -e

# Dependency on banzaicloud/cluster-registry
export GOPRIVATE='wwwin-github.cisco.com,github.cisco.com,github.com/banzaicloud'
export GONOPROXY='gopkg.in,go.uber.org'

export GOPATH=$(go env GOPATH)
export PATH="${PATH}:${GOPATH}/bin"
export GOFLAGS='-mod=readonly'

go mod download

echo "License check"
make license-cache license-check

echo "Run tests"
make test

echo "Generate"
make generate

echo "Static check"
make lint
