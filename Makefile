export PATH := $(PATH):`go env GOPATH`/bin
export GO111MODULE=on

VERSION ?= 0.68.1
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%y%m%d")
LDFLAGS := -s -w -X github.com/fatedier/frp/pkg/util/version.version=$(VERSION) -X github.com/fatedier/frp/pkg/util/version.commit=$(COMMIT) -X github.com/fatedier/frp/pkg/util/version.buildDate=$(BUILD_DATE)

.PHONY: env fmt frpc clean

all: env fmt frpc

env:
	@go version

fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

gci:
	gci write -s standard -s default -s "prefix(github.com/fatedier/frp/)" ./

vet:
	go vet -tags "frpc,noweb" ./...

frpc:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -tags "frpc,noweb" -o bin/frpc ./cmd/frpc

test: gotest

gotest:
	go test -tags "frpc,noweb" -v --cover ./cmd/...
	go test -tags "frpc,noweb" -v --cover ./client/...
	go test -tags "frpc,noweb" -v --cover ./pkg/...

clean:
	rm -f ./bin/frpc
