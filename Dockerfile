FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./cmd

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
