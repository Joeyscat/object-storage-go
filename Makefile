.PHONY: all build clean default help init test format check-license
default: help

build: build-api-server build-storage-server

build-api-server:
	go build -o build/api-server cmd/api_server/api_server.go

build-storage-server:
	go build -o build/storage-server cmd/storage_server/storage_server.go

check:
	golangci-lint run

clean:
	rm build/* -rf

# test:                   ## Run all tests.
# 	go test -v -timeout 30s ./...

init-test-env:
	bash scripts/init_test_env.sh

test-1: init-test-env
	bash scripts/test.sh