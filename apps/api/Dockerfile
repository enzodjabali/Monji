# Stage 1: Build the Go binary with CGO enabled.
FROM golang:1.18-alpine as builder

# Install build dependencies for CGO and SQLite.
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod and go.sum and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the binary (note: the main package is in cmd/api).
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

# Stage 2: Create a minimal image for running the binary.
FROM alpine:latest

# Install SQLite libraries needed at runtime.
RUN apk add --no-cache sqlite-libs

WORKDIR /root/

# Copy the built binary from the builder stage.
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]