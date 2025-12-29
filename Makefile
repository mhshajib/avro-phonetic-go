.PHONY: test lint fmt vet

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	@echo "No linter configured. Consider golangci-lint for CI."
