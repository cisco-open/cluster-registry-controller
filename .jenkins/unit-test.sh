#!/bin/bash

set -x
set -e

export GOSUMDB=off
export GOPRIVATE='github.cisco.com'
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
