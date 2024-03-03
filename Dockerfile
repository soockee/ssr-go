# Stage 1: Build the Go binary
FROM golang:latest AS builder

WORKDIR /app

# Copy the Go project files
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest && templ generate

# Build the Go binary for the desired architecture (amd64 in this case)
RUN CGO_ENABLED=0 go build -o myapp


# Stage 2: Create a minimal production image
FROM scratch

WORKDIR /app

# Copy only the binary from the previous stage
COPY --from=builder /app/myapp .

COPY ./assets ./assets

# Expose the port that your application listens on
EXPOSE 3000

# Command to run the executable
CMD ["./myapp"]
