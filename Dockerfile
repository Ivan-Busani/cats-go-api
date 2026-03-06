# FROM golang:1.26-alpine

# WORKDIR /app

# COPY go.mod ./
# RUN go mod download

# COPY . .

# EXPOSE 8080

# CMD ["go", "run", "."]

# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]