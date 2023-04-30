# Build variables
VERSION ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
UNAME := $(shell uname)
LDFLAGS += -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}

# Image URL to use all building/pushing image targets
IMG ?= controller:latest

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"
REPO_ROOT=$(shell git rev-parse --show-toplevel)
GOTESTSUM_VERSION = 0.6.0
KUBEBUILDER_TOOLS_VERSION=1.19.2
KUBEBUILDER_ASSETS_BINARY_DIR = kubebuilder-tools/${KUBEBUILDER_TOOLS_VERSION}/bin
LICENSEI_VERSION = 0.5.0
GOLANGCI_VERSION = 1.45.2
CHART_NAME = cluster-registry

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
ifeq ($(UNAME), Darwin)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_darwin_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif
ifeq ($(UNAME), Linux)
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_linux_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}
endif

.PHONY: manifests
manifests: ensure-tools ## Generate manifests
	cd api/v1alpha1 && ${REPO_ROOT}/bin/controller-gen $(CRD_OPTIONS) object:headerFile="${REPO_ROOT}/hack/boilerplate.go.txt" paths="./..." output:crd:artifacts:config=${REPO_ROOT}/deploy/charts/cluster-registry/crds

.PHONY: test
test: bin/gotestsum ensure-tools fmt vet # Run tests
	KUBEBUILDER_ATTACH_CONTROL_PLANE_OUTPUT=true KUBEBUILDER_ASSETS="$${PWD}/bin/${KUBEBUILDER_ASSETS_BINARY_DIR}" bin/gotestsum ./... -coverprofile cover.out
	cd deploy/charts && go test

.PHONY: manager
manager: fmt vet binary ## Build manager binary

.PHONY: binary
binary:
	go build -ldflags "${LDFLAGS}" -o bin/manager ./cmd/manager

.PHONY: run
run: fmt vet ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./cmd/manager/

fmt: ## Run go fmt against code
	go fmt ./...

vet: ## Run go vet against code
	go vet ./...

.PHONY: ensure-tools
ensure-tools:
	@echo "ensure tools"
	@scripts/download-deps.sh
	@scripts/install_kubebuilder_tools.sh --kubebuilder-tools-version=${KUBEBUILDER_TOOLS_VERSION}

# This will clean up the local kubebuilder-tools directory and re-download the assets again.
.PHONY: download-kubebuilder-tools
download-kubebuilder-tools:
	@rm -r bin/kubebuilder-tools/${KUBEBUILDER_TOOLS_VERSION}
	@scripts/install_kubebuilder_tools.sh --kubebuilder-tools-version=${KUBEBUILDER_TOOLS_VERSION}

# Build the docker image
docker-build:
	docker build . -t ${IMG} --build-arg GITHUB_ACCESS_TOKEN="${GITHUB_ACCESS_TOKEN}"

# Push the docker image
docker-push:
	docker push ${IMG}

bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@

.PHONY: license-check
license-check: bin/licensei ## Run license check
	bin/licensei check
	bin/licensei header

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	bin/licensei cache

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

DISABLED_LINTERS ?= --disable=gci --disable=goimports --disable=gofumpt
.PHONY: lint
lint: bin/golangci-lint ## Run linter
# "unused" linter is a memory hog, but running it separately keeps it contained (probably because of caching)
	bin/golangci-lint run --disable=unused -c .golangci.yml --timeout 2m
	bin/golangci-lint run -c .golangci.yml --timeout 2m

.PHONY: lint-fix
lint-fix: bin/golangci-lint ## Run linter & fix
# "unused" linter is a memory hog, but running it separately keeps it contained (probably because of caching)
	bin/golangci-lint run --disable=unused -c .golangci.yml --fix
	bin/golangci-lint run -c .golangci.yml --fix

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## Execute go mod tidy
	go mod tidy
	go mod download all
