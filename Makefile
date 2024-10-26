COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))

all: add-copyright tidy format test

.PHONY: format
format: # format code.
	@gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # add license to file headers.
	@addlicense -v -f $(ROOT_DIR)/LICENSE $(ROOT_DIR)

.PHONY: tidy
tidy: # Handle packkages.
	@go mod tidy

.PHONY: test
test:
	@go test ./...
