# Use Go 1.25 Alpine image
FROM golang:1.25-alpine

# Disable CGO to ensure static binary (required on Alpine)
ENV CGO_ENABLED=0

WORKDIR /app

# Copy go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN go build -o main ./cmd/app

# Expose app port
EXPOSE 8080

# Run the binary
CMD ["/app/main"]
