BIN_AIR = $(shell go env GOPATH)/bin/air
$(BIN_AIR):
	go install github.com/air-verse/air@latest

.PHONY: mac-init
mac-init:

.PHONY: linux-init
linux-init: $(BIN_AIR)

.PHONY: google-cloud-auth-init
google-cloud-auth-init:
	gcloud auth login
	gcloud auth application-default login
