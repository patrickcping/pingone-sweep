TEST?=$$(go list ./...)
SWEEP_DIR=./internal/sweep
BINARY=pingone-sweep-${NAME}
VERSION=0.1.0

default: build

tools:
	go generate -tags tools tools/tools.go

build: fmtcheck
	go install -ldflags="-X main.version=$(VERSION)"

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=all $(SWEEPARGS) -timeout 10m

vet:
	@echo "==> Running go vet..."
	@go vet ./... ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)

lint: golangci-lint importlint

golangci-lint:
	@echo "==> Checking source code with golangci-lint..."
	@golangci-lint run ./...

importlint:
	@echo "==> Checking source code with importlint..."
	@impi --local . --scheme stdThirdPartyLocal ./...

devcheck: build vet tools generate lint test sweep testacc

.PHONY: tools build generate test testacc sweep vet fmtcheck depscheck lint golangci-lint importlint
