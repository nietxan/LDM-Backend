FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./cmd

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd

CMD ["app"]

