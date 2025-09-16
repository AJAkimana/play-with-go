# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN go build -o playground main.go

# Runtime stage
FROM golang:1.24-alpine

WORKDIR /app
COPY --from=builder /app/playground .

EXPOSE 8080

CMD ["./playground"]