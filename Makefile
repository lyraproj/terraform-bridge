OS_TYPE=$(shell echo `uname`| tr '[A-Z]' '[a-z]')
ifeq ($(OS_TYPE),darwin)
	OS := osx
else
	OS := linux
endif

HAS_REQUIRED_GO := $(shell go version | grep -E 'go[2-9]|go1.1[2-9]|go1.11.[4-9]')

LDFLAGS += -X "$(PACKAGE_NAME)/pkg/version.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "$(PACKAGE_NAME)/pkg/version.BuildTag=$(shell git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///')"
LDFLAGS += -X "$(PACKAGE_NAME)/pkg/version.BuildSHA=$(shell git rev-parse --short HEAD)"

GO_PLUGINS := $(subst cmd/,,$(wildcard cmd/goplugin-*))

PHONY+= default
default: LINTFLAGS = --fast
default: everything

PHONY+= all
all: LDFLAGS += -s -w # Strip debug information
all: TESTFLAGS = --race
all: everything

PHONY+= everything
everything: check-mods clean lint test plugins sanity

PHONY+= plugins
plugins: check-mods
	@$(foreach plugin,$(GO_PLUGINS),$(call build,goplugins/$(subst goplugin-,,$(plugin)),cmd/$(plugin)/main.go);)

PHONY+= harness
harness: check-mods
	@$(call build,test-harness,cmd/test-harness/main.go)

PHONY+= test
test:
	@echo "🔘 Testing ... (`date '+%H:%M:%S'`)"
	@go test github.com/lyraproj/terraform-bridge/...

PHONY+= clean
clean:
	@echo "🔘 Cleaning ..."
	@rm -rf build

PHONY+= lint
lint: $(GOPATH)/bin/golangci-lint
	@$(foreach plugin,$(GO_PLUGINS),$(call checklint,cmd/$(plugin)/...);)

PHONY+= sanity
sanity: harness
	@$(foreach plugin,$(GO_PLUGINS),$(call checksanity,$(subst goplugin-,,$(plugin)));)

PHONY+= generate
generate:
	@echo "🔘 Regenerating bridge plugins (tf-gen) ... (`date '+%H:%M:%S'`)"
	@go run cmd/tf-gen/main.go
	@echo "✅ Generation complete (`date '+%H:%M:%S'`)"
	@echo "🔘 Rebuilding ... (`date '+%H:%M:%S'`)"
	@$(MAKE) plugins

PHONY+= check-mods
check-mods:
	@echo "🔘 Ensuring go version is 1.11.4 or later (`date '+%H:%M:%S'`)"
	@if [ "$(HAS_REQUIRED_GO)" = "" ]; \
	then \
		echo "🔴 must be running Go version 1.11.4 or later.  Please upgrade and run go clean -modcache"; \
		exit 1; \
	fi
	@echo "✅ Go version is sufficient (`date '+%H:%M:%S'`)"
	@echo "🔘 Ensuring go mod is available and turned on  (`date '+%H:%M:%S'`)"
	@GO111MODULE=on go mod download || (echo "🔴 The command 'GO111MODULE=on go mod download' did not return zero exit code (exit code was $$?)"; exit 1)
	@echo "✅ Go mod is available (`date '+%H:%M:%S'`)"

define build
	echo "🔘 Building - $(1) (`date '+%H:%M:%S'`)"
	mkdir -p build/
	GO111MODULE=on go build -ldflags '$(LDFLAGS)' -o build/$(1) $(2)
	echo "✅ Build complete - $(1) (`date '+%H:%M:%S'`)"
endef

define checklint
	echo "🔘 Linting $(1) (`date '+%H:%M:%S'`)"
	lint=`GO111MODULE=on golangci-lint run $(LINTFLAGS) $(1)`; \
	if [ "$$lint" != "" ]; \
	then echo "🔴 Lint found"; echo "$$lint"; exit 1;\
	else echo "✅ Lint-free (`date '+%H:%M:%S'`)"; \
	fi
endef

define checksanity
	echo "🔘 Sanity checking $(1) (`date '+%H:%M:%S'`)"
	sanity=`GO111MODULE=on build/test-harness build/goplugins/$(1)`; \
	if [ "$$?" != 0 ]; \
	then echo "🔴 Sanity check failed"; echo "$$sanity"; exit 1;\
	else echo "✅ Sanity check complete (`date '+%H:%M:%S'`)"; \
	fi
endef

$(GOPATH)/bin/golangci-lint:
	@echo "🔘 Installing golangci-lint... (`date '+%H:%M:%S'`)"
	@GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: $(PHONY)
