BIN="./bin"
SRC=$(shell find . -name "*.go")
BINARY_NAME=search

default: all

all: fmt test

fmt:
	$(info ******************** checking format ********************)
	@test -z $(shell gofmt -w $(SRC)) || (gofmt -d $(SRC); exit 1)

build:
	$(info ******************** building ********************)
	go build -o "${BINARY_NAME}" ./cmd/search

buildall:
	$(info ******************** building all ********************)
	GOARCH=amd64 GOOS=darwin go build -o "${BINARY_NAME}-darwin" ./cmd/search
	GOARCH=amd64 GOOS=linux go build -o "${BINARY_NAME}-linux" ./cmd/search
	GOARCH=amd64 GOOS=windows go build -o "${BINARY_NAME}-windows" ./cmd/search

clean:
	go clean
	rm -rf $(BIN)
	rm -f .cp.out

test:
	$(info ******************** running tests ********************)
	go test -v -race ./...	

cover:
	$(info ******************** running tests w coverage ********************)
	go test -v -race ./... -coverprofile .cp.out
