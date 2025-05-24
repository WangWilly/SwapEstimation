################################################################################
# Stage 1: Build the application
FROM golang:1.24-alpine AS BUILD
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/main cmd/main.go

################################################################################
# Stage 2: Runtime container
FROM alpine:latest AS RUNNER
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=BUILD /app/bin/main /app/bin/main
COPY public /app/

# Command to run the executable
CMD ["/app/bin/main"]
