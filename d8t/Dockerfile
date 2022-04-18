FROM golang:alpine AS builder

RUN apk add --no-cache ca-certificates git build-base

WORKDIR /app
COPY .. .

RUN go build -o gostarter main.go

CMD tail -f /dev/null

FROM alpine:latest as release
#
LABEL maintainer="support@win.com"

WORKDIR /app

COPY --from=builder /app/gostarter .
COPY --from=builder /app/config.yaml .

CMD ["./gostarter"]