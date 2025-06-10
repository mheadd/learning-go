FROM golang:1.23-alpine

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy go mod and sum files first
COPY go.mod ./
COPY go.sum ./

# Copy config file
COPY config.json ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main .

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

EXPOSE 8080

# Switch to non-root user
USER appuser

CMD ["./main"]