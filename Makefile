GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all install-deps vet unit-test sequential-test test fmt lint cov module-updates install help

all: install-deps vet test

install-deps:
	go mod download
	go mod verify

vet:
	go vet github.com/piprate/splash/...

unit-test: vet
	go test github.com/piprate/splash/...

test: unit-test

fmt:
	gofmt -s -w ${GOFILES_NOVENDOR}

lint:
	golangci-lint run

cov:
	go test github.com/piprate/splash/... -coverprofile coverage.out;go tool cover -html=coverage.out

module-updates:
	go list -u -m -json all | go-mod-outdated -direct -update

help:
	@echo ''
	@echo ' Targets:'
	@echo '--------------------------------------------------'
	@echo ' all              - Run everything                '
	@echo ' fmt              - Format code                   '
	@echo ' vet              - Run vet                       '
	@echo ' test             - Run all tests                 '
	@echo ' unit-test        - Run unit tests                '
	@echo ' lint             - Run golangci-lint             '
	@echo '--------------------------------------------------'
	@echo ''
