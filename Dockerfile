# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies for CGO (required for sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o auth-server main.go

# Stage 2: Final minimal image
FROM alpine:latest

# Use a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser

# Copy only the binary and templates
COPY --from=builder /app/auth-server .
COPY --from=builder /app/templates ./templates

# Expose port and define environment variables
ENV PORT=8080
ENV DB_PATH=/home/appuser/data/auth.db
EXPOSE 8080

# Create volume for persistent storage
VOLUME ["/home/appuser/data"]

CMD ["./auth-server"]
