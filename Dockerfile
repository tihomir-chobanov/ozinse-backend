# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code of the project
COPY . .

# Build the application tightly for Linux architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api/main.go

# Stage 2: Final lightweight execution environment
FROM alpine:latest  

WORKDIR /app

# Install ca-certificates in case your app makes external HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the swagger docs folder since main.go needs to read it at runtime
COPY --from=builder /app/cmd/api/docs ./cmd/api/docs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["./main"]