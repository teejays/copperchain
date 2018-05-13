GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=copperchain

clean:
	rm $(BINARY_NAME)
build:
	$(GOBUILD) -o $(BINARY_NAME)
run-example:
	cd example && $(GOBUILD) -o example && ./example
