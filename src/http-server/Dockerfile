FROM golang:1.15 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/http-server/cmd

ENTRYPOINT ["/app"]