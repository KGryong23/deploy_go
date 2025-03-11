# Sử dụng base image Golang
FROM golang:1.23 AS builder  

WORKDIR /app
COPY . .

RUN go mod tidy && go mod download
RUN go build -o server

FROM alpine:latest  
WORKDIR /root/

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
