GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=copperchain.out
EXAMPLE_BINARY_NAME=example.out
clean:
	rm $(BINARY_NAME)
build:
	$(GOBUILD) -o $(BINARY_NAME)
run-example:
	cd example && $(GOBUILD) -o $(EXAMPLE_BINARY_NAME) && ./$(EXAMPLE_BINARY_NAME)
