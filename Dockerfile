# Use Go 1.25 Alpine image
FROM golang:1.25-alpine

# Install Air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Expose app port for Gin
EXPOSE 8080

# Run Air by default
CMD ["air", "-c", ".air.toml"]
