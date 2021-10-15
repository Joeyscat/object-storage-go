.PHONY: all build clean default help init test format check-license
default: help

build: build-api-server build-storage-server

build-api-server:
	go build -o build/api-server cmd/api_server/api_server.go

build-storage-server:
	go build -o build/storage-server cmd/storage_server/storage_server.go

build-oshell:
	go build -o build/oshell cmd/oshell/main.go

mockgen:
	mockgen --build_flags=--mod=mod -self_package=github.com/joeyscat/object-storage-go/internal/api_server/store -destination internal/api_server/store/mock_store.go -package store github.com/joeyscat/object-storage-go/internal/api_server/store Factory,BucketStore,ObjectStore
	mockgen --build_flags=--mod=mod -self_package=github.com/joeyscat/object-storage-go/internal/api_server/service/v1 -destination internal/api_server/service/v1/mock_service.go -package v1 github.com/joeyscat/object-storage-go/internal/api_server/service/v1 Service,BucketSrv,ObjectSrv

lint:
	golangci-lint run

clean:
	rm build/* -rf

# test:                   ## Run all tests.
# 	go test -v -timeout 30s ./...

init-test-env:
	bash scripts/init_test_env.sh

debug: init-test-env
	bash scripts/test.sh