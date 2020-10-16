# @echo off
.SILENT:

# Default repository
REPO="github.com/bitmaelum/bitmaelum-suite"

# Our defined apps and tools
APPS := bm-server bm-client bm-config
TOOLS := hash-address jwt proof-of-work readmail update-resolver resolve vault-edit resolve-auth update-pow jwt-validate migrate-hash check-org toaster

CROSS_APPS := $(foreach app,$(APPS),cross-$(app))
CROSS_TOOLS := $(foreach tool,$(TOOLS),cross-$(tool))

# Platforms we can build on for cross platform. Should be in <os>-<arch> notation
PLATFORMS := windows-amd64 linux-amd64 darwin-amd64
BUILD_ALL_PLATFORMS := $(foreach platform,$(PLATFORMS),build-all-$(platform))


# Generate LD flags
PKG=$(shell go list ./internal)
BUILD_DATE=$(shell date)
COMMIT=$(shell git rev-parse HEAD)
LD_FLAGS = -ldflags "-X '${PKG}.buildDate=${BUILD_DATE}' -X '${PKG}.gitCommit=${COMMIT}'"

# Set environment variables from GO env if not set explicitly already
ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif
ifndef $(GOOS)
    GOOS=$(shell go env GOOS)
    export GOOS
endif
ifndef $(GOARCH)
    GOARCH=$(shell go env GOARCH)
    export GOARCH
endif

# paths to binaries
GO_STATCHECK_BIN = $(GOPATH)/bin/staticcheck
GO_INEFF_BIN = $(GOPATH)/bin/ineffassign
GO_GOCYCLO_BIN = $(GOPATH)/bin/gocyclo
GO_GOIMPORTS_BIN = $(GOPATH)/bin/goimports
GO_LINT_BIN = $(GOPATH)/bin/golint

# ---------------------------------------------------------------------------

# Downloads external tools as it's not available by default
$(GO_TEST_BIN):
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/gordonklaus/ineffassign
	go get -u github.com/fzipp/gocyclo
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint

lint:
	$(GO_GOIMPORTS_BIN) -w  --format-only .

## Runs all tests for the whole repository
test: $(GO_TEST_BIN) test_goimports test_vet test_golint test_staticcheck test_ineffassign test_gocyclo test_unit

test_goimports:
	echo "Check goimports"
	$(GO_GOIMPORTS_BIN) -l .

test_vet:
	echo "Check vet"
	go vet ./...

test_staticcheck:
	echo "Check static"
	$(GO_STATCHECK_BIN) ./...

test_golint:
	echo "Check lint"
	$(GO_LINT_BIN) ./...

test_ineffassign:
	echo "Check ineffassign"
	$(GO_INEFF_BIN) ./*

test_gocyclo:
	echo "Check gocyclo"
	$(GO_GOCYCLO_BIN) -over 15 .

test_unit:
	echo "Check unit tests"
	go test ./...

clean: ## Clean releases
	go clean

# Build default OS/ARCH apps in root release directory
$(APPS):
	$(info -   Building app $@)
	go build $(LD_FLAGS) -o release/$@ $(REPO)/cmd/$@

# Build default OS/ARCH tools in root release directory
$(TOOLS):
	$(info -   Building tool $@)
	go build $(LD_FLAGS) -o release/$@ $(REPO)/tools/$@

# Build GOOS/GOARCH apps in separate release directory
$(CROSS_APPS):
	$(info -   Building app $(subst cross-,,$@) (${GOOS}-${GOARCH}))
	go build $(LD_FLAGS) -o release/${GOOS}-${GOARCH}/$(subst cross-,,$@) $(REPO)/cmd/$(subst cross-,,$@)

# Build GOOS/GOARCH tools in separate release directory
$(CROSS_TOOLS):
	$(info -   Building tool $(subst cross-,,$@) (${GOOS}-${GOARCH}))
	go build $(LD_FLAGS) -o release/${GOOS}-${GOARCH}/$(subst cross-,,$@) $(REPO)/tools/$(subst cross-,,$@)

$(BUILD_ALL_PLATFORMS): $(CROSS_APPS) $(CROSS_TOOLS)

$(PLATFORMS):
	$(eval GOOS=$(firstword $(subst -, ,$@)))
	$(eval GOARCH=$(lastword $(subst -, ,$@)))
	$(info - Cross platform build $(GOOS) / $(GOARCH))
	make -j build-all-$(GOOS)-$(GOARCH)

info:
	$(info Building BitMaelum apps and tools)

cross-info:
	$(info Cross building BitMaelum apps and tools)

build-all: cross-info $(PLATFORMS) ## Build all cross-platform binaries

build: info $(APPS) $(TOOLS) ## Build default platform binaries

all: test build ## Run tests and build default platform binaries

help: ## Display available commands
	echo "BitMaelum make commands"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
