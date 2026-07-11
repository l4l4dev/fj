.PHONY: fmt check-fmt vet test build install uninstall verify pre-commit

INSTALL_DIR ?= $(HOME)/.local/bin
INSTALL_PATH := $(INSTALL_DIR)/fj

fmt:
	go fmt ./...

check-fmt:
	@test -z "$$(gofmt -l .)" || { echo "Go files need formatting:" >&2; gofmt -l . >&2; exit 1; }

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

install:
	mkdir -p "$(INSTALL_DIR)"
	go build -o "$(INSTALL_PATH)" ./cmd/fj

uninstall:
	rm -f "$(INSTALL_PATH)"

verify: check-fmt
	git diff --check
	$(MAKE) vet
	$(MAKE) test

pre-commit: verify
