# Makefile for VM Management API Project

# Define Go binary name
BINARY_NAME=vm-api
BUILD_DIR=build

version:
	@echo 0.0.1

tidy:
	@go mod tidy

update:
	@go get -u -t ./...

tools:
	@go install go.uber.org/mock/mockgen@latest

deps: tools
	@go mod download

mock: deps
	@go generate ./...

test: mock
	@go test -json -gcflags=-l ./... ./... -covermode=atomic

# Run static analysis with go vet
vet:tidy update
	@echo "Running go vet..."
	go vet ./...

# Build the binary and place it in the build folder
build: vet
	@echo "Building the binary..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .