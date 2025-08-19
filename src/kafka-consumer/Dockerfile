FROM golang:1.25 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/kafka-consumer/cmd

ENTRYPOINT ["/app"]
