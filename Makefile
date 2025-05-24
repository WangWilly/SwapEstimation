# Makefile for swap-estimation

# Go parameters
BINARY_NAME=swap-estimation
MAIN_PATH=cmd/main.go
BUILD_DIR=bin
DOCKER_IMAGE_NAME=swap-estimation-app
DOCKER_TAG=latest

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: all build clean test run deps docker docker-migration local-db-setup help

all: deps test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	./scripts/test.sh

run:
	./scripts/dev.sh

dev: run

deps:
	$(GOMOD) download

docker:
	./scripts/build.sh

# Help target
help:
	@echo "Make commands for $(BINARY_NAME):"
	@echo "  build             - Build the binary"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  run/dev           - Run the development server"
	@echo "  deps              - Download dependencies"
	@echo "  docker            - Build Docker image for the app"
	@echo "  all               - Run deps, test and build"
