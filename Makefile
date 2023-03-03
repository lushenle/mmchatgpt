
# Image URL to use all building/pushing image targets
IMG ?= ishenle/mmchatgpt:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes
# Options are set to exit when a recipe line exits non-zero or a piped command fails
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: build
build: fmt vet ## Build binary
	go build -installsuffix cgo -ldflags="-s -w" -o bin/mmchatgpt main.go

.PHONY: run
run: fmt vet ## Run a from your host
	go run ./main.go

.PHONY: docker-build
docker-build:  ## Build docker image with the mmchatgpt
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the mmchatgpt
	docker push ${IMG}
