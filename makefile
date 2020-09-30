TARGETNAME = velocity-limit-app

ifeq ($(OS), Windows_NT)
	TARGET = $(TARGETNAME).exe
else
	TARGET = $(TARGETNAME)
endif


ifndef GO_OS
GO_OS=windows
endif

ifndef GO_ARCH
GO_ARCH=amd64
endif

ifndef CGOENABLED
CGOENABLED=1
endif

ifndef BINARY_PATH
BINARY_PATH=./bin/$(GO_OS)_$(GO_ARCH)/$(TARGETNAME)
endif

LINTERCOMMAND=golangci-lint
## GolangCI-Lint version
GOLANGCI_VERSION=1.31.0
GOLANGCI_COMMIT=fb74c2e8e99afd50cef720595ccad516160c3974

COVERAGE_THRESHOLD=80

packages = ./validate \
           ./io
# global command
.PHONY: all
all: dependencies build test coverage code-quality-print

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor


.PHONY: build
build:
	echo "**Building linux binary**"
	GOOS=linux GOARCH=${GO_ARCH} go build -o ./bin/linux_$(GO_ARCH)/$(TARGETNAME) ./cmd/velocitylimit
	cp ./cmd/velocitylimit/config.json ./bin/linux_$(GO_ARCH)/config.json
	cp ./cmd/velocitylimit/input.txt ./bin/linux_$(GO_ARCH)/input.txt
	
	echo "**Building windows exe**"
	GOOS=windows GOARCH=${GO_ARCH} go build -o ./bin/windows_$(GO_ARCH)/$(TARGETNAME).exe ./cmd/velocitylimit
	cp ./cmd/velocitylimit/config.json ./bin/windows_$(GO_ARCH)/config.json
	cp ./cmd/velocitylimit/input.txt ./bin/windows_$(GO_ARCH)/input.txt


.PHONY: test
test:
	echo "Unit testing"
	@$(foreach package,$(packages), \
		set -e; \
		go test -coverprofile $(package)/cover.out -covermode=count $(package);)

PHONY: cover-func
coverage:
	echo "mode: count" > cover-all.out
	@$(foreach package,$(packages), \
		tail -n +2 $(package)/cover.out >> cover-all.out;)
	go tool cover -func=cover-all.out


.PHONY: code-quality-print 
code-quality-print:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.31.0
	./bin/golangci-lint --exclude-use-default=false --out-format tab run ./...

