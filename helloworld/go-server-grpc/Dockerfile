FROM golang:1.20 as builder
WORKDIR /go/src
COPY .  .

RUN CGO_ENABLED=0 go build -o /go/bin/app ./main.go

FROM alpine:3.15

WORKDIR /app

COPY --from=builder /go/bin/app app
COPY manifest/config/config.toml manifest/config/config.toml

CMD /app/app
