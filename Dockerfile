FROM golang:1.23 as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app ./cmd

EXPOSE 8888

CMD ["./app"]