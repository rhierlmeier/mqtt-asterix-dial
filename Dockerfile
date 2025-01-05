# Use the official Golang image to create a build artifact.
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o mqtt-dial cmd/mqtt-dial/main.go

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/mqtt-dial .

# Command to run the executable
CMD ["./mqtt-dial"]