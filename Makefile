# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------

# CI execution flag
IS_CI_EXECUTION 	?= 0

# Build flags
GIT_COMMIT 			:= $(shell git rev-parse HEAD)
GIT_DIRTY			?= $(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE 			?= $(shell TZ=UTC0 git show --quiet --date='format-local:%Y-%m-%dT%T%z' --format="%cd")

# Test scope
TEST_PACKAGE ?= all


# -----------------------------------------------------------------------------
# System Values
# -----------------------------------------------------------------------------

# Project root directory
PROJECT_DIR := $(shell pwd)

# Project modules (relative paths)
MODULES := $(shell find . -name go.mod -exec dirname {} \; )
MODULE_NAMES := $(shell find . -name go.mod -exec dirname {} \; | sed 's|^\./||' | tr '\n' ' ' | sed 's/,$$//')

TEST_GIT_COMMIT := test123
TEST_GO_VERSION := go1.2.3
TEST_BUILD_DATE := 0000-00-00T00:00:00+0000

#LDFLAGS := -X "main.BuildDate=$(BUILD_DATE)" \
#           -X "main.GitCommit=$(GIT_COMMIT)" \
#           -s -w # -s to omit symbol table, -w to omit DWARF debugging info

# -----------------------------------------------------------------------------
# Recipes
# -----------------------------------------------------------------------------

default: build

# Format with gofmt
.PHONY: format
format:
	@printf "Running gofmt... "
	@gofmt -l -w ${PROJECT_DIR}/. 2> /dev/null
	@echo "Done!"

# Run go mod tidy on all modules
.PHONY: tidy
tidy:
	@echo "Tidying all modules..."
	@$(foreach mod,$(MODULE_NAMES), \
		printf "  - Tidying \"$(mod)\"... "; \
		cd $(PROJECT_DIR)/$(mod); \
		go mod tidy; \
		[ $$? -eq 0 ] && printf "Success\n" || printf "FAILED\n";)
ifeq ($(shell echo $$?), 0)
	@echo "Done!"
endif

# Build modules
.PHONY: build
build: #format
	@echo "Building modules..."
	@$(foreach mod,$(MODULE_NAMES), \
		printf "  - Building \"$(mod)\"... "; \
		go build -ldflags="-X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(mod).GitCommit=$(GIT_COMMIT)' -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(mod).BUILD_DATE=$(BUILD_DATE)'" ./$(mod) 2> /dev/null; \
		[ $$? -eq 0 ] && printf "Success\n" || printf "FAILED\n";)
ifeq ($(shell echo $$?), 0)
	@echo "Done!"
endif

# Initialize and populate Go workspace
.PHONY: work
work:
	@echo "Initializing Go workspace..."
	@go work init
	@$(foreach mod,$(MODULE_NAMES), \
		echo "Adding module: $(mod)"; \
		go work use $(mod);)
	@echo ""
	@echo "Done!"

# Check for missing copyright headers
.PHONY: copyright-check
copyright-check:
	@echo "Checking for missing file headers..."
	@copywrite headers --config .copywrite.hcl --plan

# Add copywrite headers to all files
.PHONY: copyright
copyright:
	@echo "Adding any missing file headers..."
	@copywrite headers --config .copywrite.hcl

# Scan modules with gosec
.PHONY: sec
sec:
	@echo "Running gosec..."
	@gosec -quiet ./...
	@echo ""
	@echo "gosec check passed!"

# Run all tests
.PHONY: test
test: test-unit test-acc

# Run unit tests
.PHONY: test-unit
test-unit:
ifeq ($(TEST_PACKAGE), all)
	@$(foreach mod,$(MODULE_NAMES), \
		echo "Running $(mod) unit tests..."; \
		echo "---"; \
		go test -v -race -ldflags="-s -w -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(mod).GoVersion=$(TEST_GO_VERSION)' -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(mod).GitCommit=$(TEST_GIT_COMMIT)' -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(mod).BuildDate=$(TEST_BUILD_DATE)'" ./$(mod);)
else
	@echo "Running ${TEST_PACKAGE} unit tests..."
	@echo "---"
	@go test -v -race -ldflags="-s -w -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(TEST_PACKAGE).GoVersion=$(TEST_GO_VERSION)' -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(TEST_PACKAGE).GitCommit=$(TEST_GIT_COMMIT)' -X 'github.com/PaloAltoNetworks/cortex-cloud-go/$(TEST_PACKAGE).BuildDate=$(TEST_BUILD_DATE)'" ./$(TEST_PACKAGE)
endif

# Run acceptance tests
.PHONY: test-acc
test-acc: build
	@echo "Running acceptance tests..."
	@TF_ACC=1 go test -v -cover -race $$(go list ./... | grep /acceptance)


# TODO: doc generation
# TODO: copywrite
