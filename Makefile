# Default config
VERSION   ?= $(shell echo `git describe --tag 2>/dev/null || git rev-parse --short HEAD` | sed -E 's|^v||g')
VERBOSE   ?= false
DISTFOLDER = dist

# Go configuration
GOBIN     := $(shell which go)
GOENV     ?=
GOOPT     ?=
GOLDF      = -X main.Version=$(VERSION)

# FPM configuration
FPMBIN    := $(shell which fpm)
FPMFLAGS  ?=

# Docker configuration
DOCKBIN   := $(shell which docker)
DOCKIMG   := blippar/mgo-exporter
DOCKOPTS  += --build-arg VERSION="$(VERSION)"

# If run as 'make VERBOSE=true', it will pass the '-v' option to GOBIN and will restore docker build output
ifeq ($(VERBOSE),true)
GOOPT     += -v
FPMFLAGS  += --verbose
else
DOCKOPTS  += -q
.SILENT:
endif

# Binary targets configuration
EXPORTERBIN= bin/mgo-exporter
TARGETS   := $(EXPORTERBIN)
GOPKGDIR   = $(@:bin/%=./cmd/%)

# Package targets configuration
DEBPKG     = $(DISTFOLDER)/mgo-exporter_$(VERSION).amd64.deb
RPMPKG     = $(DISTFOLDER)/mgo-exporter_$(VERSION).x86_64.rpm
FPMPKGS    = $(DEBPKG) $(RPMPKG)

# Create FPMFLAGS and FPMFILES from config
FPMFLAGS  += -n "mgo-exporter" -v $(VERSION) --force \
             --config-files /etc/sysconfig/mgo-exporter \
             --post-install packager/postinst.sh --post-uninstall packager/postuninst.sh
FPMFILES  += $(EXPORTERBIN)=/usr/bin/ \
             packager/sysconfig/mgo-exporter=/etc/sysconfig/ \
             packager/systemd/mgo-exporter.service=/usr/lib/systemd/system/

# Local meta targets
all: $(TARGETS)
exporter: $(EXPORTERBIN)

# Build binaries with GOBIN using target name & GOPKGDIR
$(TARGET): GOOPT += -ldflags '$(GOLDF)'
$(TARGETS):
	$(info >>> Building $@ from $(GOPKGDIR) using $(GOBIN))
	$(if $(GOENV),$(info >>> with $(GOENV) and GOOPT=$(GOOPT)),)
	$(GOENV) $(GOBIN) build -o $@ $(GOPKGDIR) $(GOOPT)

# Build binaries staticly
static: GOLDF += -extldflags "-static"
static: GOENV += CGO_ENABLED=0 GOOS=linux
static: $(TARGETS)

# Run tests using GOBIN
test: GOPKGLIST = $(shell $(GOBIN) list ./... | grep -v vendor)
test: GOOPT += -ldflags '$(GOLDF)'
test:
	$(info >>> Testing ./... using $(GOBIN))
	$(GOENV) $(GOBIN) test $(GOOPT) -cover $(GOPKGLIST)

# Packaging
rpm: FPMFLAGS += -t rpm -s dir -p $(DISTFOLDER)/NAME_VERSION.ARCH.rpm -a x86_64 --rpm-os linux
rpm: $(RPMPKG)
deb: FPMFLAGS += -t deb -s dir -p $(DISTFOLDER)/NAME_VERSION.ARCH.deb -a x86_64
deb: $(DEBPKG)

$(FPMPKGS): static
	$(info >>> Building package $@ using $(FPMBIN))
	mkdir -p $(DISTFOLDER)
	$(FPMBIN) $(FPMFLAGS) $(FPMFILES)

# Docker
docker:
	$(info >>> Building docker image $(DOCKIMG) using $(DOCKBIN))
	$(DOCKBIN) build $(DOCKOPTS) -t $(DOCKIMG):$(VERSION) -t $(DOCKIMG):latest .

# Distribuables
dist: DOCKOPTS += --no-cache
dist: rpm deb docker

# JSON-Schema generation
generate_schemas: ./scripts/generate_schemas.go
	$(info >>> Generating JSON-Schema from Go sources using $?)
	$(GOBIN) run $?

# Clean
clean:
	$(info >>> Cleaning up binaries and distribuables)
	rm -rv $(FPMPKGS) $(TARGETS)

# Always execute these targets
.PHONY: all $(TARGETS) $(FPMPKGS)
.PHONY: exporter test
.PHONY: static rpm deb docker
