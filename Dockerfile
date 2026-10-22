FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
ENV GOTOOLCHAIN=auto
RUN go mod download
COPY . .
RUN go build -o server ./day06

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]