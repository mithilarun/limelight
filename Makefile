.PHONY: build fmt test clean all

BINARY_NAME=limelight
GO=go
GOFMT=gofmt

all: fmt test build

fmt:
	@echo "Running gofmt..."
	@$(GOFMT) -w .

test:
	@echo "Running tests..."
	@$(GO) test ./...

build: fmt test
	@echo "Building binary..."
	@$(GO) build -o $(BINARY_NAME) ./cmd/limelight

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@$(GO) clean
