# Alpine Linux with golang 1.22.Alpine is small and fast
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go files to the working directory
COPY . .

# Build the Go app
RUN go build -o agent ./cmd/agent/main.go

# Start app
CMD ["./agent"]