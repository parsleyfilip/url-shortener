FROM golang:1.17-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener .

# Use a smaller image for the final stage
FROM alpine:latest

# Add ca-certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/url-shortener .

# Expose port
EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV PORT=8080
ENV MONGODB_URI=mongodb://localhost:27017
ENV MONGODB_DATABASE=url_shortener

# Run the binary
CMD ["./url-shortener"] 