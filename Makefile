.PHONY: all

SRC_PKGS=$(shell go list ./... | grep -v "vendor")
UNIT_TEST_PACKAGES=$(shell go list ./... | grep -vE 'mocks')

fmt:
	go fmt $(SRC_PKGS)

test: 
	go test ${UNIT_TEST_PACKAGES}