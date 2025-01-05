# This Makefile intended to be POSIX-compliant (2024 edition).
#
# More info: <https://pubs.opengroup.org/onlinepubs/9799919799/utilities/make.html>
.POSIX:
.SUFFIXES:


#
# PUBLIC MACROS
#

CLI     = fwdbm
CLIDIR  = ./cmd/fwdbm
DESTDIR = ./dist
GO      = go
GOFLAGS = 
LDFLAGS = -ldflags "-s -w -X main.AppVersion=$(CLI_VERSION)"

BETTERALIGN = $$($(GO) env GOPATH)/bin/betteralign
ERRCHECK = $$($(GO) env GOPATH)/bin/errcheck
STATICCHECK = $$($(GO) env GOPATH)/bin/staticcheck


#
# INTERNAL MACROS
#

CLI_CURRENT_VER_TAG   = $$(git tag --points-at HEAD | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_LATEST_VERSION    = $$(git tag | grep "^cli" | sed 's/^cli\/v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)
CLI_PSEUDOVERSION     = $$(VER="$(CLI_LATEST_VERSION)"; echo "$${VER:-2025.01.05}")-$$(TZ=UTC git --no-pager show --quiet --abbrev=12 --date='format-local:%Y%m%d%H%M%S' --format='%cd-%h')
CLI_VERSION           = $$(VER="$(CLI_CURRENT_VER_TAG)"; echo "$${VER:-$(CLI_PSEUDOVERSION)}")
MODULE_LATEST_VERSION = $$(git tag | grep "^v" | sed 's/^v//' | sort -t. -k 1,1n -k 2,2n -k 3,3n | tail -1)


#
# DEVELOPMENT TASKS
#

.PHONY: all
all: install-dependencies

.PHONY: clean
clean:
	@echo '# Delete build directories' >&2
	rm -rf $(DESTDIR)

.PHONY: info
info:
	@printf '# OS info: '
	@uname -rsv;
	@echo '# Development dependencies:'
	@$(GO) version || true
	@$(BETTERALIGN) -V=full || true
	# @$(ERRCHECK) --version || true  # not supported, see: https://github.com/kisielk/errcheck/issues/254
	@$(STATICCHECK) --version || true
	@echo '# Go environment variables:'
	@$(GO) env || true

.PHONY: check
check:
	@echo '# Unit tests' >&2
	$(GO) test -race -vet=off ./...
	@echo '# Static analysis' >&2
	$(GO) vet ./...
	$(STATICCHECK) ./...
	$(ERRCHECK) ./...
	$(GO) mod verify
	@echo '# Formatting' >&2
	$(GO) fmt ./...
	$(BETTERALIGN) ./...
	$(GO) mod tidy

.PHONY: build
build:
	@echo '# Build CLI executable: $(DESTDIR)/$(CLI)' >&2
	$(GO) build -C $(CLIDIR) $(GOFLAGS) $(LDFLAGS) -o '../../$(DESTDIR)/$(CLI)'
	@echo '# Add executable checksum to: $(DESTDIR)/$(CLI).sha256sum' >&2
	cd $(DESTDIR); sha256sum $(CLI) >> $(CLI).sha256sum

.PHONY: dist
dist: clean fwdbm-linux_amd64 fwdbm-windows_amd64.exe

.PHONY: install-dependencies
install-dependencies:
	@echo '# Install development dependencies:' >&2
	$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@latest
	$(GO) install github.com/kisielk/errcheck@latest
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest
	@echo '# Install CLI dependencies' >&2
	@GOFLAGS='-v -x' $(GO) get -C $(CLIDIR) $(GOFLAGS) .

.PHONY: cli-release
cli-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new CLI release tag' >&2
	@VER="$(CLI_LATEST_VERSION)"; printf 'Choose new version number for CLI (calver; >%s): ' "$${VER:-2025.01.05}"
	@read -r NEW_VERSION; \
		git tag "cli/v$$NEW_VERSION"; \
		git push --tags

.PHONY: module-release
module-release: check test
	@echo '# Update local branch' >&2
	@git pull --rebase
	@echo '# Create new Go module release tag' >&2
	@VER="$(MODULE_LATEST_VERSION)"; printf 'Choose new version number for module (semver; >%s): ' "$${VER:-0.0.0}"
	@read -r NEW_VERSION; \
		git tag "v$$NEW_VERSION"; \
		git push --tags


#
# SUPPORTED EXECUTABLES
#

# this force using `go build` to changes detection in Go project (instead of `make`)
.PHONY: fwdbm-linux_amd64 \
 fwdbm-windows_amd64.exe

fwdbm-linux_amd64:
	@GOOS=linux GOARCH=amd64 $(MAKE) CLI=$@ build

fwdbm-windows_amd64.exe:
	@GOOS=windows GOARCH=amd64 $(MAKE) CLI=$@ build
