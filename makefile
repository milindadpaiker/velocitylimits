TARGETNAME = platform-extended-attributes-retrieval-service

ifeq ($(OS), Windows_NT)
	TARGET = $(TARGETNAME).exe
else
	TARGET = $(TARGETNAME)
endif

ifndef BINARY_PATH
BINARY_PATH=$(GOPATH)/bin/linux_amd64/$(TARGETNAME)
endif

ifndef GO_OS
GO_OS=linux
endif

ifndef GO_ARCH
GO_ARCH=amd64
endif

LINTERCOMMAND=golangci-lint
## GolangCI-Lint version
GOLANGCI_VERSION=1.31.0
GOLANGCI_COMMIT=fb74c2e8e99afd50cef720595ccad516160c3974

COVERAGE_THRESHOLD=80

.PHONY: code-quality-print ## Run golang-cilint with printing to stdout
code-quality-print: bin/golangci-lint
	./bin/golangci-lint --exclude-use-default=false --out-format tab run ./...

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin 1.31.0
