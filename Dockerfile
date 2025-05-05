# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

# --- Runtime Stage ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY config ./config

EXPOSE 8080

ENTRYPOINT ["./main"]
