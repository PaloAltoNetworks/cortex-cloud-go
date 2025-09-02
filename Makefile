# Get project root directory to allow for execution from other
# filesystem locations (e.g. "make -C /path/to/cortex-cloud-go build")
PROJECT_DIR=$(shell pwd)

# Define the Go modules to be included in the workspace.
# These are the relative paths to the root of each module.
MODULES := \
    . \
    ./api/ \
    ./appsec/ \
    ./cloudonboarding/ \
    ./enums/ \
    ./errors/ \
    ./internal/app/ \
    ./log/ \
    ./xsiam/ \


default: build

# Format with gofmt
.PHONY: format
format:
	@echo "Running gofmt..."
	@gofmt -l -w ${PROJECT_DIR}/.

# Build provider binary
.PHONY: build
build: format
	@echo "Building SDK plugin..."
	@go build ${PROJECT_DIR}
ifeq ($(shell echo $$?), 0)
	@echo "Done!"
endif

# Init go workspace / Uses MODULES
.PHONY: init
init:
	@echo "Initializing Go workspace..."
	@rm -f go.work go.work.sum
	go work init
	@for mod in $(MODULES); do \
		echo "Adding module: $$mod"; \
		go work use $$mod; \
	done
	@echo "Workspace initialized successfully."

# Test current directory and all subdirs
.PHONY: test
test:
	@echo "Running tests for all workspace modules..."
	@go test -v ./...
	@echo "Testing complete"

# TODO: acceptance tests
# TODO: doc generation
# TODO: copywrite
