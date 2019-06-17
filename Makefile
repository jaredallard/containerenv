TARGETS           ?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x
DIST_DIRS         = find * -type d -exec
APP               = containerenv

# For testing
CLUSTER := pontus-integration-testing
PROJECT := azuqua-218321
ZONE := us-west1-b

# go option
GO        ?= go
PKG       := $(shell glide novendor)
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin
BINARIES  := containerenv
SHORT_NAME := containerenv

# Required for globs to work correctly
SHELL=/bin/bash


.PHONY: all
all: build

.PHONY: release
release:
	@./hack/make-release.sh

.PHONY: dep
dep:
	@go mod vendor

.PHONY: lint
lint:
	./hack/verify-golint.sh

.PHONY: gofmt
gofmt:
	gofmt -w -s pkg/
	gofmt -w -s cmd/

.PHONY: build
build: lint
	@echo "  building releases in ./bin/..."
	$(GO) build $(GOFLAGS) -o "$(BINDIR)/containerenv" -tags '$(TAGS)' -ldflags '$(LDFLAGS)' github.com/jaredallard/containerenv/cmd/

.PHONY: clean
clean:
	@rm -rf $(BINDIR)

HAS_GOX := $(shell command -v gox;)
HAS_GIT := $(shell command -v git;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GOX
	go get -u github.com/mitchellh/gox
endif

ifndef HAS_GIT
	$(error You must install Git)
endif

include versioning.mk