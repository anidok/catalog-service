.PHONY: all

SRC_PKGS=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell go list ./... | grep -vE 'mocks')

fmt:
	go fmt $(SRC_PKGS)

test: 
	go test ${UNIT_TEST_PACKAGES}

lint:
	golangci-lint cache clean
	golangci-lint run --config .golangci.yml -v

setup-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.2
	go install github.com/vektra/mockery/v2@latest

generate-mocks:
	mockery --name=Client --dir=internal/opensearch --output=test/mocks/opensearch --outpkg=opensearch
	mockery --name=ServiceRepository --dir=internal/repository --output=test/mocks/repository --outpkg=repository
	mockery --name=ServiceUsecase --dir=internal/usecase --output=test/mocks/usecase --outpkg=usecase

migrate:
	curl -X DELETE "http://localhost:9200/services"
	go run cmd/migrate/main.go

ingest:
	go run cmd/ingest/main.go

migrate-ingest: migrate ingest
all: fmt test lint setup-tools generate-mocks migrate-ingest
