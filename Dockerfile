FROM golang:1.21-alpine as builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /kode .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /kode .
COPY application.yaml .

EXPOSE 80

CMD ["./kode"]
