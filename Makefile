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