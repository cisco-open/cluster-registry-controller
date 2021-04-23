#!/bin/bash

set -x
set -e

# get artifactory creds
sed -i '1d' "${NYOTA_CREDENTIALS_FILE}"
. "${NYOTA_CREDENTIALS_FILE}"

# Dependency on banzaicloud/cluster-registry
export GOPRIVATE='github.cisco.com,github.com/banzaicloud'
export GOPROXY="https://proxy.golang.org,https://${artifactory_user}:${artifactory_password}@${artifactory_url},direct"

export GOPATH=$(go env GOPATH)
export PATH="${PATH}:${GOPATH}/bin"
export GOFLAGS='-mod=readonly'

go mod download

echo "License check"
#make license-check

echo "Run tests"
make test

echo "Generate"
make generate

echo "Static check"
make lint
