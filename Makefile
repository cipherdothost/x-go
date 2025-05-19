# SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
#
# SPDX-License-Identifier: MIT

.POSIX:
.SUFFIXES:

GO    ?= go
GIT   ?= git
REUSE ?= reuse
RM    ?= rm

REQUIRED_TOOLS := $(GO) $(GIT)
GO_MIN_VERSION := 1.24

all: pre-commit

check-tools: # Checks that all required tools for working with this repository are installed.
	$(foreach tool,$(REQUIRED_TOOLS),\
		$(if $(shell command -v $(tool)),,$(error "$(tool) not found in PATH")))
	@$(GO) version | awk 'NR==1 {if ($$3 < "go$(GO_MIN_VERSION)") exit 1}'

pre-commit: check-tools tidy fmt lint vulnerabilities test # Runs all pre-commit checks.

commit: pre-commit # Commits the changes to the repository.
	$(GIT) commit -s

doc: check-tools # Serves the documentation locally.
	$(GO) run golang.org/x/tools/cmd/godoc@latest -http=localhost:1967

tidy: check-tools # Updates the go.mod file to ensure it matches the source code.
	$(GO) mod tidy

fmt: check-tools # Formats Go source files in this repository.
	$(GO) run mvdan.cc/gofumpt@latest -e -extra -w .

lint: check-tools # Runs linters on the codebase.
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

lint/licenses: # Run reuse to ensure the project complies with the REUSE specification.
	$(REUSE) lint --lines

vulnerabilities: check-tools # Analyzes the codebase and looks for vulnerabilities affecting it.
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

test: check-tools # Runs unit tests.
	$(GO) test -cover -race -vet all -mod readonly ./...

test/coverage: check-tools # Generates a coverage profile and open it in a browser.
	$(GO) test -coverprofile cover.out ./...
	$(GO) tool cover -html=cover.out

licenses: check-tools # Runs go-licenses to check the licenses of the dependencies and generate a CSV file.
	$(GO) run github.com/google/go-licenses@latest report \
		--template '.github/license-3rdparty.tpl' \
		--ignore 'go.cipher.host/x' \
		./... > LICENSE-3rdparty.csv

.PHONY: all check-tools pre-commit commit doc tidy fmt lint \
	lint/licenses vulnerabilities test test/coverage licenses
