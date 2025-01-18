GO_SOURCES=$(shell find . -name "*.go")

BIN_AIR = $(shell go env GOPATH)/bin/air
$(BIN_AIR):
	go install github.com/air-verse/air@latest

BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint
$(BIN_GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

dist/prd/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/dev/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/loc/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/e2e/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ go/cmd/$(notdir $@)/*.go

.PHONY: server-init-loc
server-init-loc: $(BIN_AIR)
	$(BIN_AIR) -c ./air.server.loc.toml

.PHONY: test
test: dist/e2e/server
	# E2E
	mkdir -p cov/e2e/server && rm -f cov/e2e/server/*
	FILE_PATH_SERVER_BIN=$(abspath dist/e2e/server) \
	GOCOVERDIR=$(abspath cov/e2e/server) \
	go test -v -count=1 ./e2e/...
	# GOCOVERDIR=cov/e2e/server sh ../../run-e2e.sh dist/e2e/server "sh ../../wait-until-http-health.sh http://localhost:8080/health" "go test -v -count=1 ./e2e/..."
	# UT
	mkdir -p cov/go && rm -f cov/go/*
	go test -v -cover ./go/... -args -test.gocoverdir=$(abspath cov/go)

.PHONY: merge-test-report
merge-test-report:
	go tool covdata percent -i=cov/e2e/server,cov/go -o=textfmt.0.txt
	go tool cover -html=textfmt.0.txt -o=gocov.html
	go tool cover -func=textfmt.0.txt -o=gocovfunc.txt

.PHONY: lint
lint: $(BIN_GOLANGCI_LINT)
	$(BIN_GOLANGCI_LINT) run ./...
