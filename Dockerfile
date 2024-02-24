# Use the official Go image as a base
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module and Go sum files
COPY go.mod ./

# Copy the rest of the application source code
COPY . .

# Build the application
# Build the Go application
RUN go build -o /app/bin/main cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/app/bin/main"]