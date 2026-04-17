.PHONY: build install clean

BINARY_NAME=kodia

# Build the binary
build:
	go build -o bin/$(BINARY_NAME) cmd/kodia/main.go

# Install the binary to GOPATH/bin
install:
	go install cmd/kodia/main.go

# Remove build artifacts
clean:
	rm -rf bin/
