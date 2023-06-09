FROM golang:1.19-buster as builder

WORKDIR /soa-hmw-2
COPY . .

WORKDIR /soa-hmw-2/cmd/client
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o .

FROM alpine:3.15.4
WORKDIR /soa-hmw-2
COPY --from=builder /soa-hmw-2/cmd/client .
CMD ["./client"]
