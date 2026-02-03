# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o scout9 ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS requests to GRID API
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/scout9 .

# Create data directory for JSON backup
RUN mkdir -p data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./scout9"]
