# Step 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the dependencies (and cache them)
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pingpong .

# Step 2: Create the final image
FROM alpine:latest

# Install required dependencies (e.g., for HTTP and Prometheus handling)
RUN apk --no-cache add ca-certificates

# Set the working directory for the application
WORKDIR /root/

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/pingpong .

# Copy the environment file (if necessary) or any additional resources
COPY .env ./

# Expose the necessary ports (Service port and Metrics port)
EXPOSE 8080 9080

# Set the entry point to run the application
CMD ["./pingpong"]