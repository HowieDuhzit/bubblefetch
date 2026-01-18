.PHONY: build install clean run test

# Build the binary
build:
	go build -o bubblefetch ./cmd/bubblefetch

# Build with optimizations
build-release:
	go build -ldflags="-s -w" -o bubblefetch ./cmd/bubblefetch

# Install to system
install: build-release
	sudo mv bubblefetch /usr/local/bin/

# Clean build artifacts
clean:
	rm -f bubblefetch

# Run the program
run: build
	./bubblefetch

# Run with a specific theme
run-dracula: build
	./bubblefetch --theme dracula

run-minimal: build
	./bubblefetch --theme minimal

run-nord: build
	./bubblefetch --theme nord

run-gruvbox: build
	./bubblefetch --theme gruvbox

run-tokyo: build
	./bubblefetch --theme tokyo-night

run-monokai: build
	./bubblefetch --theme monokai

run-solarized: build
	./bubblefetch --theme solarized-dark

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install system-wide
install: build-release
	./install.sh

# Uninstall
uninstall:
	./uninstall.sh

# Run benchmark
bench: build
	./bubblefetch --benchmark

# Export examples
export-json: build
	./bubblefetch --export json --pretty=true

export-yaml: build
	./bubblefetch --export yaml

export-text: build
	./bubblefetch --export text

# Build example plugin
plugin-hello:
	go build -buildmode=plugin -o plugins/hello.so plugins/examples/hello.go

# Build all example plugins
plugins: plugin-hello

# Install plugins to user directory
install-plugins: plugins
	mkdir -p ~/.config/bubblefetch/plugins
	cp plugins/*.so ~/.config/bubblefetch/plugins/

# Clean plugin artifacts
clean-plugins:
	rm -f plugins/*.so
