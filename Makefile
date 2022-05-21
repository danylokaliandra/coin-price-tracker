GO      := go
BIN     := tracklet
BINPATH := /usr/local/bin


all: tidy fmt lint test build

tidy:
	$(info ▶ cleaning dependencies...)
	$(GO) mod tidy

fmt:
	$(info ▶ formatting...)
	gofmt -s -w .

lint:
	$(info ▶ running lint tools...)
	golangci-lint run -v

test:
	$(info ▶ running tests...)
	go test -coverprofile cover.out -v ./...

build:
	$(info ▶ compiling program...)
	$(GO) build -o $(BIN)

install: build
	$(info ▶ installing program...)
	mkdir -p ~/.tracklet
	cp ./config/tracklet.yaml ~/.tracklet/tracklet.yaml
	sudo cp $(BIN) $(BINPATH)/$(BIN)

uninstall:
	$(info ▶ uninstalling program...)
	cp ~/.tracklet/tracklet.yaml /tmp/tracklet_backup.yaml
	rm -rf ~/.tracklet
	sudo rm $(BINPATH)/$(BIN)

clean:
	$(info ▶ removing binary...)
	rm $(BIN)


.PHONY: tidy fmt lint test build install uninstall clean
