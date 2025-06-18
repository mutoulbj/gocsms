# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gocsms ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gocsms .
COPY .env.example .env
EXPOSE 3000 9000
CMD ["./gocsms"]