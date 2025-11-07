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

# Module to test (defaults to all)
UNIT_TEST_MODULE 				?= all
# Unit test to run
UNIT_TEST_NAME				?= ""

#------------------------------------------------------------------------------
# LDFLAGS (Linker Flags) Definitions
#------------------------------------------------------------------------------

define LDFLAGS_template
-s -w \
-X '$(1).GitCommit=$(GIT_COMMIT)' \
-X '$(1).CortexServerVersion=$(CORTEX_SERVER_VERSION)' \
-X '$(1).CortexPAPIVersion=$(CORTEX_PAPI_VERSION)' \
-X '$(1).BuildDate=$(BUILD_DATE)' \
-X '$(1).GoVersion=$(GO_VERSION)'
endef

define TEST_LDFLAGS_template
-s -w \
-X '$(1).GitCommit=$(TEST_GIT_COMMIT)' \
-X '$(1).CortexServerVersion=$(TEST_CORTEX_SERVER_VERSION)' \
-X '$(1).CortexPAPIVersion=$(TEST_CORTEX_PAPI_VERSION)' \
-X '$(1).BuildDate=$(TEST_BUILD_DATE)' \
-X '$(1).GoVersion=$(TEST_GO_VERSION)'
endef

#------------------------------------------------------------------------------
# Phony Targets
#------------------------------------------------------------------------------

.PHONY: help format tidy lint build work work-sync copyright-check copyright sec test test-unit test-acc tag clean

#------------------------------------------------------------------------------
# Main Targets
#------------------------------------------------------------------------------

help: ## Show this help message.
	@sed -rn 's/^([^:]+):.*[ ]##[ ](.+)/\1:\2/p' $(MAKEFILE_LIST) | column -ts:

format: ## Format all Go source files.
	@echo "Running gofmt..."
	@gofmt -l -w .
	@echo "Done."

# DO NOT CHANGE THIS ORDER
tidy: ## Tidy all go.mod files.
	@echo "Tidying all modules..."
	@echo "  - log"
	@(cd ./log && rm -f go.sum && go mod tidy)
	@echo "  - errors"
	@(cd ./errors && rm -f go.sum && go mod tidy)
	@echo "  - enums"
	@(cd ./enums && rm -f go.sum && go mod tidy)
	@echo "  - types"
	@(cd ./types && rm -f go.sum && go mod tidy)
	@echo "  - internal/config"
	@(cd ./internal/config && rm -f go.sum && go mod tidy)
	@echo "  - internal/client"
	@(cd ./internal/client && rm -f go.sum && go mod tidy)
	@echo "  - appsec"
	@(cd ./appsec && rm -f go.sum && go mod tidy)
	@echo "  - cloudonboarding"
	@(cd ./cloudonboarding && rm -f go.sum && go mod tidy)
	@echo "  - cwp"
	@(cd ./cwp && rm -f go.sum && go mod tidy)
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

work-sync: ## Synchronize dependencies defined in Go workspace file (go.work)
	@echo "Syncronizing Go workspace dependencies..."
	@go work sync
	@echo "Syncronization successful."

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

test: test-unit test-acc # Run unit and acceptance tests.

# To run the unit tests for the platform module:
#
# UNIT_TEST_MODULE=platform make test-unit
#
# To run only the "TestClient_GetUserGroup" unit test:
#
# UNIT_TEST_NAME=TestClient_GetUserGroup make test-unit
test-unit: ## Run unit tests. Can be scoped to a specific module and/or test.
ifeq ($(UNIT_TEST_MODULE), all)
	@if [ -z "$(UNIT_TEST_NAME)" ]; then \
		echo "Running unit tests for all modules..."; \
	else \
		echo "Running unit test \"$(UNIT_TEST_NAME)\" for all modules..."; \
	fi
	@$(foreach mod,$(MODULE_NAMES), \
		echo "--- Running tests for $(mod) ---"; \
		go test -v -race $(if $(UNIT_TEST_NAME),-run $(UNIT_TEST_NAME),) -ldflags="$(call TEST_LDFLAGS_template,$(mod))" ./$(mod) || exit 1; \
	)
	@echo "All unit tests passed."
else
	@if [ -z "$(UNIT_TEST_NAME)" ]; then \
		echo "--- Running tests for $(UNIT_TEST_MODULE) ---"; \
	else \
		echo "--- Running test \"$(UNIT_TEST_NAME)\" for $(UNIT_TEST_MODULE) ---"; \
	fi
	@go test -v -race $(if $(UNIT_TEST_NAME),-run $(UNIT_TEST_NAME),) -ldflags="$(call TEST_LDFLAGS_template,$(UNIT_TEST_MODULE))" ./$(UNIT_TEST_MODULE)
endif

test-acc: build ## Run acceptance tests.
	@echo "Running acceptance tests..."
	@TF_ACC=1 go test -v -cover -race $$(go list ./... | grep /acceptance)

clean: ## Clean up workspace files.
	@echo "Cleaning workspace..."
	@-rm -f go.work go.work.sum
	@echo "Done."
