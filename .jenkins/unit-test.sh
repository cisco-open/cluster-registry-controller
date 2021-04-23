#!/bin/bash

set -x
set -e

export GOSUMDB=off
# Dependency on banzaicloud/cluster-registry
export GOPRIVATE='github.cisco.com,github.com/banzaicloud'
export GOPROXY="https://proxy.golang.org,direct"

export GOPATH=$(go env GOPATH)
export GOFLAGS='-mod=readonly'
export PATH="${PATH}:${GOPATH}/bin"

echo "License check"
#make license-check

echo "Run tests"
make test

echo "Generate"
make generate

echo "Static check"
make lint
