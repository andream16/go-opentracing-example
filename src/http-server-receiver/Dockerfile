FROM golang:1.25 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/http-server-receiver/cmd

ENTRYPOINT ["/app"]
