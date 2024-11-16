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
RUN go build -o /main ./cmd

# Base image for the MongoDB service
FROM mongo:6.0 AS mongodb

# Expose MongoDB default port
EXPOSE 27017

# Base image for Redis
FROM redis:8.0 AS redis

# Expose Redis default port
EXPOSE 6379

# Final stage to run all services
FROM ubuntu:22.04

# Copy the Go binary from the build stage
COPY --from=go-builder /main /usr/local/bin/main

# Install necessary tools to run services
RUN apt-get update && apt-get install -y \
    supervisor && rm -rf /var/lib/apt/lists/*

# Create supervisor configuration file
RUN mkdir -p /etc/supervisor/conf.d
