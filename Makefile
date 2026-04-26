TEMPLCMD=templ
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=main

# Tailwind CSS
TAILWIND=./tailwindcss
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),arm64)
TAILWIND_ARCH=macos-arm64
else
TAILWIND_ARCH=macos-x64
endif
TAILWIND_URL=https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$(TAILWIND_ARCH)

all: test build

$(TAILWIND):
	curl -sL $(TAILWIND_URL) -o $(TAILWIND)
	chmod +x $(TAILWIND)

css: $(TAILWIND)
	$(TAILWIND) -i assets/css/input.css -o assets/css/output.css --minify

build: css
	$(TEMPLCMD) generate
	$(GOBUILD) -o $(BINARY_NAME) -v .

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f assets/css/output.css

run: css
	$(TEMPLCMD) generate
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)

css-watch: $(TAILWIND)
	$(TAILWIND) -i assets/css/input.css -o assets/css/output.css --watch
