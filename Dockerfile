# Stage 1: Build the Go application
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Stage 2: Create a minimal image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
EXPOSE 8080
CMD ["./main"]
