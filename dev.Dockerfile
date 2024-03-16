# Dockerfile for local development

# Use a base image with Go installed
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to install dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Install gin for hot-reloading
RUN go install github.com/codegangsta/gin@latest

# Copy the source code into the container
COPY . .

# Expose port 3000 for gin
EXPOSE 8100

# Command to run the application with gin
CMD ["gin", "run", "-p", "8100", "cmd/main.go"]
