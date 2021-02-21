FROM golang:1.16 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/http-server-initiator/cmd

ENTRYPOINT ["/app"]