.PHONY:
.SILENT:
.DEFAULT_GOAL := run

VERSION ?= $(strip $(shell ./scripts/version.sh))
VERSION_NUMBER := $(strip $(shell ./scripts/version.sh number))
COMMIT_HASH := $(shell git rev-parse --short HEAD)
GO_OPT_BASE := -ldflags "-X main.version=$(VERSION) $(GO_LDFLAGS) -X main.commitHash=$(COMMIT_HASH)"
OUT_BIN ?= ./.bin/dv-updater

OUT_DIR ?= ./.bin
## Build:


build:
	go build $(GO_OPT_BASE) -o $(OUT_BIN) ./cmd/app

build-linux:
		GOOS=linux GOARCH=amd64 go build $(GO_OPT_BASE) -o $(OUT_BIN)-linux ./cmd/app; \

run: build
	$(OUT_BIN) $(filter-out $@,$(MAKECMDGOALS))

run-linux:
	$(OUT_BIN)-linux $(filter-out $@,$(MAKECMDGOALS))

## Lint:
lint:
	golangci-lint run --show-stats

fmt:
	gofumpt -l -w .