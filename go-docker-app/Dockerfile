# Stage 1: Build the Go application
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application binary for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory inside the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 9090

# Run the compiled Go binary
CMD ["./main"]
