# Use the official Golang image as a build stage
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Run tests
RUN go test ./... -v

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o vigor-api ./cmd/vigor-api

# Use a minimal image to run the binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/vigor-api .

# Command to run the executable
CMD ["/app/vigor-api"]
