# Base image for Golang
FROM golang:1.23.3 AS go-builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd

# Final stage
FROM ubuntu:22.04

# Copy the Go binary from the build stage
COPY --from=go-builder /main /usr/local/bin/main

# Install necessary tools to run services
RUN apt-get update && apt-get install -y \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

# Create supervisor configuration directory and config file
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Make the binary executable
RUN chmod +x /usr/local/bin/main

# Command to run supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
