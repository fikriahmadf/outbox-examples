.PHONY: all wire swagger run build tidy clean

all: wire swagger build

# Generate wire injectors
wire:
	cd internal/app && go run github.com/google/wire/cmd/wire

# Generate Swagger docs (requires `swag` installed)
swagger:
	swag init -g main.go -o ./docs

# Run the application
run:
	go run .

# Build the application
build:
	go build -o bin/outbox-examples .

# Go mod tidy
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin docs
