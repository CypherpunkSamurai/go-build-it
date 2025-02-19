# ------------------
#  Golang Makefile
# ------------------

# Go parameters
GOCMD=go
GOTIDY=$(GOCMD) mod tidy
GOBUILD=$(GOCMD) build
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main

# Platform-specific settings
ifeq ($(OS),Windows_NT)
    BINARY_SUFFIX=.exe
    RM=del /Q
else
    BINARY_SUFFIX=
    RM=rm -f
endif
# Add target to install binary
BINARY=$(BINARY_NAME)$(BINARY_SUFFIX)

all: build

build:
	$(GOBUILD) -o $(BINARY) -v

clean:
	$(GOCLEAN)

run:
	$(GOCLEAN)
	$(GOFMT) .
	$(GOTIDY)
	$(GOBUILD) -o $(BINARY) -v
	./$(BINARY)

format:
	$(GOFMT) .

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) ./...

.PHONY: all build clean run test deps