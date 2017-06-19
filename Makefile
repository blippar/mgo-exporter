# Default config
GOBIN	  := $(shell which go)
VERBOSE   ?= false
GOENV	  ?=

# If run as 'make VERBOSE=true', it will pass th '-v' option to GOBIN
ifeq ($(VERBOSE),true)
GOOPT     += -v
endif

# Targets configuration
EXPORTER_BIN = bin/exporter

# List all target to create a rul that manage all of them
TARGETS   := $(EXPORTER_BIN)

# Precreate a variable to get package name from binary name
PKGDIR     = $(@:bin/%=./cmd/%)

# Local meta targets
all: $(TARGETS)
exporter: $(EXPORTER_BIN)

# Check if GOBIN exists before running a rule
_check_gobin:
	$(if $(wildcard $(GOBIN)),,$(error GOBIN is not set, is go installed))

# Build binaries with GOBIN using target name & PKGDIR
$(TARGETS): _check_gobin
	$(info >>> Building $@ from $(PKGDIR) using $(GOBIN))
	$(if $(GOENV),$(info >>> with $(GOENV) and GOOPT=$(GOOPT)),)
ifeq ($(VERBOSE),true)
	$(GOENV) $(GOBIN) build $(GOOPT) -o $@ $(PKGDIR)
else
	@$(GOENV) $(GOBIN) build $(GOOPT) -o $@ $(PKGDIR)
endif

# Build binaries staticly
static: GOOPT += -ldflags '-extldflags "-static"'
static: GOENV += CGO_ENABLED=0 GOOS=linux
static: $(TARGETS)

# Run tests using GOBIN
test: _check_gobin
	$(info >>> Testing ./... using $(GOBIN))
	@$(GOBIN) test $(GOOPT) ./...

# Always execute these targets
.PHONY: all $(TARGETS)
.PHONY: exporter static test
.PHONY: _check_gobin
