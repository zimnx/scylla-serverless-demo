all: build

SHELL :=/bin/bash -euEo pipefail -O inherit_errexit

comma :=,

MAKE_REQUIRED_MIN_VERSION:=4.2 # for SHELLSTATUS
GO_REQUIRED_MIN_VERSION ?=1.18

GIT_TAG ?=$(shell git describe --long --tags --abbrev=7 --match 'v[0-9]*')$(if $(filter $(.SHELLSTATUS),0),,$(error git describe failed))
GIT_TAG_SHORT ?=$(shell git describe --tags --abbrev=7 --match 'v[0-9]*')$(if $(filter $(.SHELLSTATUS),0),,$(error git describe failed))
GIT_COMMIT ?=$(shell git rev-parse --short "HEAD^{commit}" 2>/dev/null)$(if $(filter $(.SHELLSTATUS),0),,$(error git rev-parse failed))
GIT_TREE_STATE ?=$(shell ( ( [ ! -d ".git/" ] || git diff --quiet ) && echo 'clean' ) || echo 'dirty')

GO ?=go
GO_MODULE ?=$(shell $(GO) list -m)$(if $(filter $(.SHELLSTATUS),0),,$(error failed to list go module name))
GOPATH ?=$(shell $(GO) env GOPATH)
GOOS ?=$(shell $(GO) env GOOS)
GOEXE ?=$(shell $(GO) env GOEXE)
GOFMT ?=gofmt
GOFMT_FLAGS ?=-s -l

GO_VERSION :=$(shell $(GO) version | sed -E -e 's/.*go([0-9]+.[0-9]+.[0-9]+).*/\1/')
GO_PACKAGE ?=$(shell $(GO) list -m -f '{{ .Path }}' || echo 'no_package_detected')
GO_PACKAGES ?=./...

go_packages_dirs :=$(shell $(GO) list -f '{{ .Dir }}' $(GO_PACKAGES) || echo 'no_package_dir_detected')
GO_TEST_PACKAGES ?=$(GO_PACKAGES)
GO_BUILD_PACKAGES ?=./cmd/...
GO_BUILD_PACKAGES_EXPANDED ?=$(shell $(GO) list $(GO_BUILD_PACKAGES))
go_build_binaries =$(notdir $(GO_BUILD_PACKAGES_EXPANDED))
GO_BUILD_FLAGS ?=-trimpath
GO_BUILD_BINDIR ?=
GO_LD_EXTRA_FLAGS ?=
GO_TEST_PACKAGES :=./pkg/... ./cmd/...
GO_TEST_FLAGS ?=-race
GO_TEST_COUNT ?=
GO_TEST_EXTRA_FLAGS ?=
GO_TEST_ARGS ?=
GO_TEST_EXTRA_ARGS ?=
GO_TEST_E2E_EXTRA_ARGS ?=

JQ ?=jq
YQ ?=yq -e

define version-ldflags
-X $(1).versionFromGit="$(GIT_TAG)" \
-X $(1).commitFromGit="$(GIT_COMMIT)" \
-X $(1).gitTreeState="$(GIT_TREE_STATE)" \
-X $(1).buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"
endef
GO_LD_FLAGS ?=-ldflags '$(strip $(call version-ldflags,$(GO_PACKAGE)/pkg/version) $(GO_LD_EXTRA_FLAGS))'

# TODO: look into how to make these local to the targets
export DOCKER_BUILDKIT :=1
export GOVERSION :=$(shell go version)
export KUBEBUILDER_ASSETS :=$(GOPATH)/bin
export PATH :=$(GOPATH)/bin:$(PATH):

# $1 - required version
# $2 - current version
define is_equal_or_higher_version
$(strip $(filter $(2),$(firstword $(shell printf '%s\n%s' '$(1)' '$(2)' | sort -V -r -b))))
endef

# $1 - program name
# $2 - required version variable name
# $3 - current version string
define require_minimal_version
$(if $($(2)),\
$(if $(strip $(call is_equal_or_higher_version,$($(2)),$(3))),,$(error `$(1)` is required with minimal version "$($(2))", detected version "$(3)". You can override this check by using `make $(2):=`)),\
)
endef

ifneq "$(MAKE_REQUIRED_MIN_VERSION)" ""
$(call require_minimal_version,make,MAKE_REQUIRED_MIN_VERSION,$(MAKE_VERSION))
endif

ifneq "$(GO_REQUIRED_MIN_VERSION)" ""
$(call require_minimal_version,$(GO),GO_REQUIRED_MIN_VERSION,$(GO_VERSION))
endif

# $1 - package name
define build-package
	$(if $(GO_BUILD_BINDIR),mkdir -p '$(GO_BUILD_BINDIR)',)
	$(strip CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) $(GO_LD_FLAGS) \
		$(if $(GO_BUILD_BINDIR),-o '$(GO_BUILD_BINDIR)/$(notdir $(1))$(GOEXE)',) \
	$(1))

endef

# We need to build each package separately so go build creates appropriate binaries
build:
	$(if $(strip $(GO_BUILD_PACKAGES_EXPANDED)),,$(error no packages to build: GO_BUILD_PACKAGES_EXPANDED var is empty))
	$(foreach package,$(GO_BUILD_PACKAGES_EXPANDED),$(call build-package,$(package)))
.PHONY: build

clean:
	$(RM) $(go_build_binaries)
.PHONY: clean

help:
	$(info The following make targets are available:)
	@$(MAKE) -f $(firstword $(MAKEFILE_LIST)) --print-data-base --question no-such-target 2>&1 | grep -v 'no-such-target' | \
	grep -v -e '^no-such-target' -e '^makefile' | \
	awk '/^[^.%][-A-Za-z0-9_]*:/	{ print substr($$1, 1, length($$1)-1) }' | sort -u
.PHONY: help
