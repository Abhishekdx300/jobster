FROM golang:1.24-alpine AS builder


# Set working directory
WORKDIR /app

# Copy go mod files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main


# Expose the dev port
EXPOSE 8080

CMD ["/app/main"]
