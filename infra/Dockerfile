# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Copy go.sum if it exists (optional for fresh builds)
COPY go.su[m] ./

# Download dependencies and generate go.sum if needed
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o twitter-api .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /build/twitter-api .

# Copy static files
COPY --from=builder /build/public ./public

# Expose port
EXPOSE 8080

# Run the application
CMD ["./twitter-api"]
