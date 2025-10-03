SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rule

.DEFAULT_GOAL := help

#------------------------------------------------------------------------------
# System & Build Configuration
#------------------------------------------------------------------------------

# Linker values
GIT_COMMIT 					:= $(shell git rev-parse HEAD)
CORTEX_SERVER_VERSION 		:= master-platform-v4.2.0-4877-g4886d-7fe3
CORTEX_PAPI_VERSION 		:= 1.2
BUILD_DATE 					?= $(shell TZ=UTC0 git show --quiet --date='format-local:%Y-%m-%dT%T%z' --format="%cd")
GO_VERSION 					:= $(shell go version)
# Project root directory	
PROJECT_DIR 				:= $(CURDIR)
# Project module directory paths
MODULE_PATHS 					:= $(shell find . -name go.mod -exec dirname {} \;)
# List of module names (e.g., "api", "appsec")
MODULE_NAMES 				:= $(patsubst ./%,%,$(MODULE_PATHS))

# Test values
TEST_GIT_COMMIT 			:= test123
TEST_CORTEX_SERVER_VERSION 	:= test-v0.0.0
TEST_CORTEX_PAPI_VERSION 	:= 0.0
TEST_GO_VERSION 			:= go1.2.3
TEST_BUILD_DATE 			:= 0000-00-00T00:00:00+0000

# Package to test (defaults to all)
TEST_PACKAGE 				?= all

#------------------------------------------------------------------------------
# LDFLAGS (Linker Flags) Definitions
#------------------------------------------------------------------------------

define LDFLAGS_template
-s -w \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GitCommit=$(GIT_COMMIT)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).CortexServerVersion=$(CORTEX_SERVER_VERSION)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).CortexPAPIVersion=$(CORTEX_PAPI_VERSION)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).BuildDate=$(BUILD_DATE)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GoVersion=$(GO_VERSION)'
endef

define TEST_LDFLAGS_template
-s -w \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GitCommit=$(TEST_GIT_COMMIT)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).CortexServerVersion=$(TEST_CORTEX_SERVER_VERSION)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).CortexPAPIVersion=$(TEST_CORTEX_PAPI_VERSION)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).BuildDate=$(TEST_BUILD_DATE)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GoVersion=$(TEST_GO_VERSION)'
endef

#------------------------------------------------------------------------------
# Phony Targets
#------------------------------------------------------------------------------

.PHONY: help format tidy lint build work copyright-check copyright sec test test-unit test-acc tag clean

#------------------------------------------------------------------------------
# Main Targets
#------------------------------------------------------------------------------

help: ## Show this help message.
	@sed -rn 's/^([^:]+):.*[ ]##[ ](.+)/\1:\2/p' $(MAKEFILE_LIST) | column -ts:

format: ## Format all Go source files.
	@echo "Running gofmt..."
	@gofmt -l -w .
	@echo "Done."

tidy: ## Tidy all go.mod files.
	@echo "Tidying all modules..."
	@$(foreach mod,$(MODULE_PATHS), \
		echo "  - Tidying $(mod)"; \
		(cd $(mod) && go mod tidy) || exit 1; \
	)
	@echo "Tidy successful."

lint: ## Lint all modules with go vet.
	@echo "Linting all modules..."
	@$(foreach mod,$(MODULE_NAMES), \
		echo "  - Linting $(mod)"; \
		(cd $(mod) && go vet ./...) || exit 1; \
	)
	@echo "Lint successful."

build: ## Build all modules.
	@echo "Building all modules..."
	@$(foreach mod,$(MODULE_NAMES), \
		echo "  - Building $(mod)"; \
		go build -ldflags="$(call LDFLAGS_template,$(mod))" ./$(mod) || exit 1; \
	)
	@echo "Build successful."

work: ## Initialize or update the Go workspace file (go.work).
	@echo "Setting up Go workspace..."
	@if [ ! -f "go.work" ]; then go work init; fi
	@go work use $(MODULE_PATHS)
	@go work sync
	@echo "Workspace ready."

copyright-check: ## Check for missing copyright headers.
	@echo "Checking for missing file headers..."
	@copywrite headers --config .copywrite.hcl --plan

copyright: ## Add missing copyright headers to all files.
	@echo "Adding any missing file headers..."
	@copywrite headers --config .copywrite.hcl

sec: ## Scan modules for security issues with gosec.
	@echo "Running gosec..."
	@gosec -quiet ./...
	@echo "gosec check passed!"

test: test-unit test-acc ## Run all tests.

test-unit: ## Run unit tests for all or a specific module (e.g., make test-unit TEST_PACKAGE=api).
ifeq ($(TEST_PACKAGE), all)
	@echo "Running unit tests for all modules..."
	@$(foreach mod,$(MODULE_NAMES), \
		echo "--- Running tests for $(mod) ---"; \
		go test -v -race -ldflags="$(call TEST_LDFLAGS_template,$(mod))" ./$(mod) || exit 1; \
	)
	@echo "All unit tests passed."
else
	@echo "--- Running tests for $(TEST_PACKAGE) ---"
	@go test -v -race -ldflags="$(call TEST_LDFLAGS_template,$(TEST_PACKAGE))" ./$(TEST_PACKAGE)
endif

test-acc: build ## Run acceptance tests.
	@echo "Running acceptance tests..."
	@TF_ACC=1 go test -v -cover -race $$(go list ./... | grep /acceptance)

clean: ## Clean up workspace files.
	@echo "Cleaning workspace..."
	@-rm -f go.work go.work.sum
	@echo "Done."
