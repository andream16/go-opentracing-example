FROM golang:1.25 as build

WORKDIR /build

COPY . .

RUN go build -o /app ./src/grpc-server/cmd

ENTRYPOINT ["/app"]
