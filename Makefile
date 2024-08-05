BUILD=go build -ldflags="-w -s"

default: build

.PHONY: build
build: clean vet
	@echo "Building wcrypt-cli..."
	cd cmd/cli/ && $(BUILD) -o ../../bin/wcrypt-cli.exe
	@echo "Building wcrypt-server..."
	cd cmd/server/ && $(BUILD) -o ../../bin/wcrypt-server.exe

.PHONY: clean
clean:
	@rm -rf bin/

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: all
all:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt-cli.exe cmd/cli/*
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt-server.exe cmd/server/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt-cli cmd/cli/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt-server cmd/server/*
