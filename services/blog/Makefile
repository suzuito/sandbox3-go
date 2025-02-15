GO_SOURCES=$(shell find . -name "*.go")

BIN_GOLANGCI_LINT = $(shell go env GOPATH)/bin/golangci-lint
$(BIN_GOLANGCI_LINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.5

BIN_GOLANG_MIGRATE = $(shell go env GOPATH)/bin/migrate
$(BIN_GOLANG_MIGRATE):
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

dist/prd/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/dev/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/loc/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -o $@ go/cmd/$(notdir $@)/*.go

dist/e2e/server: $(GO_SOURCES)
	mkdir -p $(dir $@) && go build -cover -o $@ go/cmd/$(notdir $@)/*.go

.PHONY: server-init-loc
server-init-loc:
	go tool air -c ./air.server.loc.toml

.PHONY: db-init
db-init: $(BIN_GOLANG_MIGRATE)
	PGPASSWORD=root psql -v ON_ERROR_STOP=1 -U root -h $(DB_HOST) -p 5432 -c 'CREATE DATABASE $(DB_NAME)' || true
	$(BIN_GOLANG_MIGRATE) -source file://./db/schema -database "postgres://root:root@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" drop -f
	$(BIN_GOLANG_MIGRATE) -source file://./db/schema -database "postgres://root:root@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" up

.PHONY: db-init-loc
db-init-loc:
	go run go/cmd/db-init-loc/*.go

.PHONY: test
test: dist/e2e/server
	# E2E
	mkdir -p cov/e2e/server && rm -f cov/e2e/server/*
	FILE_PATH_SERVER_BIN=$(abspath dist/e2e/server) \
	GOCOVERDIR=$(abspath cov/e2e/server) \
	go test -v -count=1 ./e2e/...
	# UT
	mkdir -p cov/go && rm -f cov/go/*
	go test -v -cover ./go/... -args -test.gocoverdir=$(abspath cov/go)

.PHONY: merge-test-report
merge-test-report:
	go tool covdata percent -i=cov/e2e/server,cov/go -o=textfmt.0.txt
	grep -v .sql.go textfmt.0.txt > textfmt.1.txt
	go tool cover -html=textfmt.1.txt -o=gocov.html
	go tool cover -func=textfmt.1.txt -o=gocovfunc.txt

.PHONY: lint
lint: $(BIN_GOLANGCI_LINT)
	$(BIN_GOLANGCI_LINT) run ./...

# make create-migrate NAME=hoge
.PHONY: create-migrate
create-migrate: $(BIN_GOLANG_MIGRATE)
	$(BIN_GOLANG_MIGRATE) create -ext sql -dir db/schema $(NAME)

.PHONY: sqlc
sqlc:
	go tool sqlc generate
