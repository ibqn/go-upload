FROM golang:1.25-alpine AS builder

RUN apk add --no-cache \
    gcc \
    g++ \
    make \
    vips-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

RUN apk add --no-cache \
    vips \
    ca-certificates

WORKDIR /app

COPY --from=builder /app/api .

RUN mkdir -p /app/file-storage

EXPOSE 8888

CMD ["./api"]
