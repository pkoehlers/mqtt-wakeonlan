# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Stage 2: Create a minimal container with the built binary
FROM scratch

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app /app

# Command to run the executable
ENTRYPOINT ["/app"]