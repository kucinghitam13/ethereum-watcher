FROM golang:1.19 as builder

RUN git config --global core.eol lf
RUN git config --global core.autocrlf input

COPY . /ethereum-watcher
WORKDIR /ethereum-watcher

RUN go build -o /build/ethereum-watcher main.go

FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates

FROM ubuntu
COPY --from=builder /build/ethereum-watcher /app/
COPY --from=builder /ethereum-watcher/files/ /app/files/
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app

EXPOSE 8080

ENTRYPOINT [ "./ethereum-watcher" ]

