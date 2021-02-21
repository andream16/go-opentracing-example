FROM golang:1.16 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/kafka-consumer/cmd

ENTRYPOINT ["/app"]