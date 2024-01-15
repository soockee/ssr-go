FROM golang:latest AS builder

WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 3000

ENTRYPOINT ["./main"]
