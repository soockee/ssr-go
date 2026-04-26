# Stage 1: Build the Go binary
FROM golang:latest AS builder

WORKDIR /app

# Cache Go module downloads
COPY go.mod go.sum ./
RUN go mod download

# Install tools (cached unless Go version changes)
RUN go install github.com/a-h/templ/cmd/templ@latest

# Download Tailwind CSS standalone CLI
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "arm64" ]; then TW_ARCH="linux-arm64"; else TW_ARCH="linux-x64"; fi && \
    curl -sL "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-${TW_ARCH}" -o /usr/local/bin/tailwindcss && \
    chmod +x /usr/local/bin/tailwindcss

# Copy source (cache busted only when code changes)
COPY . .

# Generate templ and build CSS
RUN templ generate && \
    tailwindcss -i assets/css/input.css -o assets/css/output.css --minify

# Build the Go binary
RUN CGO_ENABLED=0 go build -o myapp


# Stage 2: Create a minimal production image
FROM arm64v8/ubuntu:22.04

RUN apt update && apt install -y bash ca-certificates mime-support && rm -rf /var/cache/apt/*
WORKDIR /app

# Copy only the binary from the previous stage
COPY --from=builder /app/myapp .
COPY --from=builder /app/assets ./assets

# Expose the port that your application listens on
EXPOSE 443
EXPOSE 80

VOLUME ["/certs"]


# Command to run the executable
ENTRYPOINT ["./myapp"]
