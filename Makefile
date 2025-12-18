# Let's specify the name of the binary to build
BINARY_NAME=fako


# Let's build the application

build:
	@echo "Building Fako..."
	go build -o bin/$(BINARY_NAME) cmd/fako/main.go

# Let's run the application

run:
	go run cmd/fako/main.go

# Let's write the command to run tests

test:
	go test ./...

# Now, let's write the way we will clean the artifact

clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Format code

fmt:
	go fmt ./...

# the linter config
lint:
	golangci-lint run