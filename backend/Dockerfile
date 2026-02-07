# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the specific service
ARG SERVICE=api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/${SERVICE} ./${SERVICE}/cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
ARG SERVICE=api-gateway
COPY --from=builder /app/bin/${SERVICE} .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]
