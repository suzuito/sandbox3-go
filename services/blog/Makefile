GO_SOURCES=$(shell find . -name "*.go")

BIN_AIR = $(shell go env GOPATH)/bin/air
$(BIN_AIR):
	go install github.com/air-verse/air@latest

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

.PHONY: server-init-e2e
server-init-e2e: $(BIN_AIR)
	mkdir -p cov/e2e/server && rm -f cov/e2e/server/*
	GOCOVERDIR=cov/e2e/server $(BIN_AIR) -c ./air.server.e2e.toml

.PHONY: server-e2e
server-e2e:
	go test -v ./e2e/...

.PHONY: test
test:
	mkdir -p cov/go && rm -f cov/go/*
	go test -v -cover ./go/... -args -test.gocoverdir=$(abspath cov/go)

.PHONY: merge-test-report
merge-test-report:
	go tool covdata percent -i=cov/e2e/server,cov/go -o=textfmt.0.txt
	go tool cover -html=textfmt.0.txt -o=gocov.html
	go tool cover -func=textfmt.0.txt -o=gocovfunc.txt
