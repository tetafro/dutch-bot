FROM golang:1.21.0-alpine3.18 AS build

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

COPY . .

RUN go build -o ./bin/dutch-bot .

FROM alpine:3.18

WORKDIR /app

COPY --from=build /build/bin/dutch-bot /app/

RUN apk add --no-cache ca-certificates && \
    addgroup -S -g 5000 dutch-bot && \
    adduser -S -u 5000 -G dutch-bot dutch-bot && \
    chown -R dutch-bot:dutch-bot .

USER dutch-bot

ENTRYPOINT ["/app/dutch-bot"]
