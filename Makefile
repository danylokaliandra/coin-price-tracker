GO := go
BIN := tracklet
M	= $(shell printf "\033[34;1m▶\033[0m")


.PHONY: tidy build clean

all: tidy fmt lint test build

tidy:
	$(info $(M) cleaning dependencies...)
	$(GO) mod tidy

fmt:
	$(info $(M) checking formatting...)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info $(M) running lint tools...)
	golangci-lint run -v

test:
	$(info $(M) running tests...)
	go test -coverprofile cover.out -v ./...

build:
	$(info $(M) compiling program...)
	$(GO) build -o $(BIN)

clean:
	$(info $(M) removing binary...)
	rm $(BIN)