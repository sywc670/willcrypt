BUILD=go build -ldflags="-w -s"

default: build

.PHONY: build
build: clean vet
	@echo "Building wcrypt-cli..."
	cd cli/ && $(BUILD) -o ../wcrypt
	@echo "Building wcrypt-server..."
	cd server/ && $(BUILD) -o ../wcrypt-server

.PHONY: clean
clean:
	@rm -rf bin/ wcrypt.exe wcrypt-server.exe wcrypt wcrypt-server

.PHONY: clean-all
clean-all:
	@rm -rf bin/ wcrypt.exe wcrypt-server.exe wcrypt wcrypt-server id.txt pairs.txt priv.key

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: all
all:
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt.exe cli/*
	GOOS=windows GOARCH=amd64 $(BUILD) -o bin/wcrypt-server.exe server/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt cli/*
	GOOS=linux GOARCH=amd64 $(BUILD) -o bin/wcrypt-server server/*
