# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=simplebc
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
# test:
# 	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/urfave/cli
	$(GOGET) github.com/boltdb/bolt/...
	$(GOGET) github.com/google/go-cmp/cmp
	$(GOGET) golang.org/x/crypto/ripemd160
	$(GOGET) github.com/btcsuite/btcutil


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

