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

RUN CGO_ENABLED=1 go build -o go-upload .

FROM alpine:latest

RUN apk add --no-cache \
    vips \
    ca-certificates

WORKDIR /app

COPY --from=builder /app/go-upload .

RUN mkdir -p /app/file-storage

EXPOSE 8000

CMD ["./go-upload"]
