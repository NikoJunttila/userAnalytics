# Build stage
FROM docker.io/golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (using vendor if available is faster but for docker we usually let it download)
# However, the project has a vendor directory, so we can use that.
COPY vendor/ ./vendor/

# Copy the rest of the source code
COPY . .

# Build the application
# We use -mod=vendor to use the checked-in dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o analytics-app .

# Final stage
FROM docker.io/alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/analytics-app .

# Copy the schema directory for goose migrations on startup
COPY sql/schema/ ./sql/schema/

# Expose the port the app runs on
EXPOSE 8000

# Run the application
CMD ["./analytics-app"]
