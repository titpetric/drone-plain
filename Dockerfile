## build image (throw-away)

FROM golang:1.12-alpine AS builder

ADD . /drone-plain

WORKDIR /drone-plain

RUN apk add git && go build -o build/drone-plain ./cmd/drone-plain/main.go

## target image

FROM alpine:latest

RUN apk add --no-cache ca-certificates bash

COPY --from=builder /drone-plain/build/* /app/

WORKDIR /app

ENTRYPOINT ./drone-plain