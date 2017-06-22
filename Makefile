# Default config
VERSION   = 0.1.1
VERBOSE   ?= false

# Go configuration
GOBIN	  := $(shell which go)
GOENV	  ?=
GOOPT     ?=

# FPM configuration
FPMBIN    := $(shell which fpm)
FPMFLAGS  ?=

# If run as 'make VERBOSE=true', it will pass th '-v' option to GOBIN
ifeq ($(VERBOSE),true)
GOOPT     += -v
FPMFLAGS  += --verbose
endif

# Binary targets configuration
EXPORTER_BIN = bin/mgo-exporter
TARGETS     := $(EXPORTER_BIN)
GOPKGDIR     = $(@:bin/%=./cmd/%)

# Package targets configuration
DEBPKG	   = dist/mgo-export_$(VERSION).x86_64.deb
RPMPKG	   = dist/mgo-export_$(VERSION).x86_64.rpm
FPMPKGS    = $(DEBPKG) $(RPMPKG)

# Create FPMFLAGS and FPMFILES from config
FPMFLAGS  += -n "mgo-exporter" -v $(VERSION) --force \
	     --config-files /etc/sysconfig/mgo-exporter \
             --post-install packager/postinst.sh --post-uninstall packager/postuninst.sh
FPMFILES  += $(EXPORTER_BIN)=/usr/bin/ \
             packager/sysconfig/mgo-exporter=/etc/sysconfig/ \
	     packager/systemd/mgo-exporter.service=/usr/lib/systemd/system/

# Local meta targets
all: $(TARGETS)
exporter: $(EXPORTER_BIN)

# Build binaries with GOBIN using target name & GOPKGDIR
$(TARGETS):
	$(info >>> Building $@ from $(GOPKGDIR) using $(GOBIN))
	$(if $(GOENV),$(info >>> with $(GOENV) and GOOPT=$(GOOPT)),)
	$(GOENV) $(GOBIN) build $(GOOPT) -o $@ $(GOPKGDIR)

# Run tests using GOBIN
test:
	$(info >>> Testing ./... using $(GOBIN))
	@$(GOBIN) test $(GOOPT) ./...

# Build binaries staticly
static: GOOPT += -ldflags '-extldflags "-static"'
static: GOENV += CGO_ENABLED=0 GOOS=linux
static: $(TARGETS)

# Packaging
rpm: FPMFLAGS += -t rpm -s dir -p dist/NAME_VERSION.ARCH.rpm -a x86_64 --rpm-os linux
rpm: $(RPMPKG)
deb: FPMFLAGS += -t deb -s dir -p dist/NAME_VERSION.ARCH.deb -a x86_64
deb: $(DEBPKG)

$(FPMPKGS): static
	$(info >>> Building package $@ using fpm)
	mkdir -p dist
	$(FPMBIN) $(FPMFLAGS) $(FPMFILES)

# Always execute these targets
.PHONY: all $(TARGETS) $(FPMPKGS)
.PHONY: exporter static rpm deb test
