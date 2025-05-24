# Makefile for labs-hr-go

# Go parameters
BINARY_NAME=labs-hr-go
MAIN_PATH=cmd/main.go
BUILD_DIR=bin
DOCKER_IMAGE_NAME=labs-hr-go-app
DOCKER_TAG=latest
MIGRATION_IMAGE_NAME=labs-hr-go-migration
MIGRATION_TAG=latest

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

docker-migration:
	./scripts/migration_build.sh

local-db-setup:
	./scripts/local_setup.sh

# Help target
help:
	@echo "Make commands for $(BINARY_NAME):"
	@echo "  build             - Build the binary"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  run/dev           - Run the development server"
	@echo "  deps              - Download dependencies"
	@echo "  docker            - Build Docker image for the app"
	@echo "  docker-migration  - Build Docker image for migrations"
	@echo "  local-db-setup    - Set up local database with Docker"
	@echo "  all               - Run deps, test and build"
