NAME := schnegge
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: test
test:  ## Run the tests.
	@$(GO) test ./...

.PHONY: build
build: cmd/schnegge/main.go  ## Build a binary.
	$(GO) build -ldflags "$(LDFLAGS)" -o ${NAME} cmd/schnegge/main.go

.PHONY: cross
cross: cmd/schnegge/main.go  ## Build binaries for cross platform.
	mkdir -p pkg
	@# darwin
	@for arch in "amd64" "arm64"; do \
		GOOS=darwin GOARCH=$${arch} make build; \
		zip pkg/$(NAME)_$(VERSION)_darwin_$${arch}.zip $(NAME); \
	done;
	@# linux
	@for arch in "amd64"; do \
		GOOS=linux GOARCH=$${arch} make build; \
		zip pkg/$(NAME)_$(VERSION)_linux_$${arch}.zip $(NAME); \
	done;
	@# windows
	@for arch in "amd64"; do \
		GOOS=windows GOARCH=$${arch} make build; \
		zip pkg/$(NAME)_$(VERSION)_windows_$${arch}.zip $(NAME); \
	done;


.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'