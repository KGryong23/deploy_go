# Sử dụng base image Golang
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Biên dịch code
RUN go mod tidy
RUN go build -o server

# Image nhẹ chạy app
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/server .

# Expose port 8080
EXPOSE 8080

# Chạy ứng dụng
CMD ["./server"]
