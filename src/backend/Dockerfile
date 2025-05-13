# Use Golang image to build and run
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy the Go module manifests (go.mod, go.sum)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the Go code
COPY . .

# Expose the port your app runs on (default: 8000)
EXPOSE 8000

# Run the application directly
CMD ["go", "run", "main.go"]
