.PHONY: fmt check-fmt vet test build verify pre-commit

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

verify: check-fmt
	git diff --check
	$(MAKE) vet
	$(MAKE) test

pre-commit: verify
