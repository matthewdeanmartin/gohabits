# Variables
BINARY_NAME=gohabits
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/gohabits/main.go

CONFIG_DIR=~/.config/gohabits
CONFIG_FILE=$(CONFIG_DIR)/config.yaml



.PHONY: all build run clean init check

all: check build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_PATH) $(MAIN_PATH)

# Run directly
run:
	go run $(MAIN_PATH)

# Run tests and vet
check:
	go vet ./...
	go test ./...

# Clean build artifacts
#clean:
#	@echo "Cleaning..."
#	rm -f $(BINARY_NAME)

# INIT: Sets up local config (~/.config/habits)
# Does NOT touch Go source code.
# Creates a sample config.yaml if one does not exist.
init:
	 ./script/init.sh ~/.config/gohabits/config.yaml
