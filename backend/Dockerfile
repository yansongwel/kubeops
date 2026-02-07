# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the monolithic application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/kubeops ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/kubeops .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./kubeops"]
