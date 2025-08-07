# Build stage
FROM golang:1.24-alpine AS builder

# Instalar certificados CA
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Create static directory in final stage
RUN mkdir -p /static

# Final stage
FROM scratch

# Copiar certificados CA para HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/main /main

# Copy static files
COPY --from=builder /app/static /static

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/main"]
