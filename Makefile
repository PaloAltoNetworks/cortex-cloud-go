SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rule

.DEFAULT_GOAL := help

#------------------------------------------------------------------------------
# System & Build Configuration
#------------------------------------------------------------------------------

# Build info values
VERSION 		:= $(shell git describe --tags --always --dirty)
GIT_COMMIT 		:= $(shell git rev-parse HEAD)
BUILD_DATE 		?= $(shell TZ=UTC0 git show --quiet --date='format-local:%Y-%m-%dT%T%z' --format="%cd")
GO_VERSION          := $(shell go version)

# Project root directory
PROJECT_DIR := $(CURDIR)

# Find all go.mod files and list their directories.
MODULES := $(shell find . -name go.mod -exec dirname {} \;)
# Create a clean list of module names (e.g., "api", "appsec") from the paths.
MODULE_NAMES := $(patsubst ./%,%,$(MODULES))

# Test values
TEST_GIT_COMMIT := test123
TEST_GO_VERSION := go1.2.3
TEST_BUILD_DATE := 0000-00-00T00:00:00+0000
TEST_VERSION    := test-version

# Package to test (defaults to all)
TEST_PACKAGE ?= all

#------------------------------------------------------------------------------
# LDFLAGS (Linker Flags) Definitions
#------------------------------------------------------------------------------

define LDFLAGS_template
-s -w \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GitCommit=$(GIT_COMMIT)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).BuildDate=$(BUILD_DATE)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GoVersion=$(GO_VERSION)'
endef
#-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).Version=$(VERSION)' \

define TEST_LDFLAGS_template
-s -w \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GitCommit=$(TEST_GIT_COMMIT)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).BuildDate=$(TEST_BUILD_DATE)' \
-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).GoVersion=$(TEST_GO_VERSION)'
endef
#-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(1).Version=$(TEST_VERSION)' \

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
	@$(foreach mod,$(MODULES), \
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
	@go work use $(MODULES)
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

#tag: ## Create new patch version tags for all modules (e.g., make tag AUTO_APPROVE=true).
#	@echo "Calculating new tags for all modules..."
#	TAGS_TO_CREATE=""
#	$(for mod in $(MODULES); do \
#		tag_prefix=""; \
#		if [ "$$mod" != "." ]; then \
#			tag_prefix="$${mod#./}/" \
#		fi; \
#
#		LATEST_TAG=$(git tag -l "$${tag_prefix}v*" | sort -V | tail -n 1); \
#
#		if [ -z "$$LATEST_TAG" ]; then \
#			NEW_VERSION="v0.0.1"; \
#		else \
#			LATEST_TAG_VERSION=$${LATEST_TAG%$${tag_prefix}}v; \
#			LATEST_TAG_VERSION=$${LATEST_TAG_VERSION#v}; \
#			MAJOR=$$(echo "$$LATEST_TAG_VERSION" | cut -d. -f1); \
#			MINOR=$$(echo "$$LATEST_TAG_VERSION" | cut -d. -f2); \
#			PATCH=$$(echo "$$LATEST_TAG_VERSION" | cut -d. -f3); \
#			if ! [[ "$$PATCH" =~ ^[0-9]+$$ ]]; then \
#				echo "Error: Could not parse patch version from '$$LATEST_TAG' for module '$${mod}'." >&2; \
#				exit 1; \
#			fi; \
#			NEW_PATCH=$(PATCH + 1); \
#			NEW_VERSION="v$${MAJOR}.$${MINOR}.$${NEW_PATCH}"; \
#		fi;  \
#		NEW_TAG="$${tag_prefix}$${NEW_VERSION}"; \
#		echo "  - Proposing tag for $${mod}: $$NEW_TAG"; \
#		TAGS_TO_CREATE="$$TAGS_TO_CREATE $${NEW_TAG}"; \
#	done; \
#	if [ -z "$$TAGS_TO_CREATE" ]; then \
#		echo "No tags to create."; \
#		exit 0; \
#	fi; \
#	echo "\nThe following tags will be created:"; \
#	for tag in $$TAGS_TO_CREATE; do \
#		echo "  - $$tag" \
#	done; \
#	echo ""; \
#		if [ "$$AUTO_APPROVE" != "true" ]; then \
#			read -p "Apply these tags? (y/n) " -n 1 -r; \
#			echo ""; \
#			if [[ ! "$$REPLY" =~ ^[Yy]$$ ]]; then \
#				echo "Aborted by user."; \
#				exit 1; \
#			fi; \
#		fi; \
#			echo "Applying tags..."; \
#		for tag in $$TAGS_TO_CREATE; do \
#			#git tag "$$tag" \
#			echo "git tag \"$$tag\""; \
#			echo "Created tag: $$tag"; \
#		done; \
#		echo "---"; \
#		echo "Run 'git push --tags' to push them to the remote."; \
#	)

clean: ## Clean up workspace files.
	@echo "Cleaning workspace..."
	@-rm -f go.work go.work.sum
	@echo "Done."
