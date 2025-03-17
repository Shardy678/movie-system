FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go module files from the backend directory
COPY backend/go.mod backend/go.sum ./

RUN go mod download

# Copy the rest of the backend source files
COPY backend ./

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

RUN ls -lah main

EXPOSE 8080

CMD ["./main"]
