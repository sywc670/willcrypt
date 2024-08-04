BUILD=go build -ldflags="-w -s"

default: build

build:
	@echo "Building wcrypt-cli..."
	cd cmd/cli/ && $(BUILD) -o ../../bin/wcrypt-cli.exe
	@echo "Building wcrypt-server..."
	cd cmd/server/ && $(BUILD) -o ../../bin/wcrypt-server.exe

clean:
	@rm -rf bin/

all:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt-cli.exe cmd/cli/*
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt-server.exe cmd/server/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt-cli cmd/cli/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt-server cmd/server/*
