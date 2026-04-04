export PATH := $(PATH):`go env GOPATH`/bin
export GO111MODULE=on
LDFLAGS := -s -w

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
