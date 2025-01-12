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

.PHONY: init-server
init-server: $(BIN_AIR)
	$(BIN_AIR) -c ./air.server.toml
