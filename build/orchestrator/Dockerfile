# Alpine Linux with golang 1.22.Alpine is small and fast
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go files to the working directory
COPY . .

# Build the Go app
RUN go build -o orchestrator ./cmd/orchestrator/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Start app
CMD ["./orchestrator"]