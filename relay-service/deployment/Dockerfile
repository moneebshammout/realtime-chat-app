# Stage 1: Build the Go binary
FROM golang:1.21 AS builder


# Set the working directory in the container
WORKDIR /app

# Copy the Go project files (excluding specified directories) to the container
COPY go.mod go.sum ./

RUN go mod download  && go mod verify

COPY . .

# Build the Go binary
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o relay-service ./cmd/http

# Stage 2: Create a minimal image for running the binary
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/relay-service  .

# Expose the port the app runs on
EXPOSE 7000

# Command to run the binary
CMD ["./relay-service"]
