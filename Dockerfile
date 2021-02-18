FROM golang:1.15-alpine as builder

ENV GO111MODULE on

RUN apk --no-cache add git make gcc libc-dev upx

WORKDIR /go/src/github.com/IsmaelPereira/bot-telegram-isma/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN make build

###

FROM alpine:3.13

ENV CONFIG_DIR /local/config/

COPY --from=builder /go/src/github.com/IsmaelPereira/bot-telegram-isma/ /etc/telegram-bot/config/

WORKDIR /opt/telegram-bot-isma/bin/
COPY --from=builder /go/src/github.com/IsmaelPereira/bot-telegram-isma/ /opt/telegram-bot-isma/bin/telegram-bot

CMD ["/opt/telegram-bot-isma/bin/telegram-bot","/etc/telegram-bot/config/"]

RUN apk --no-cache add ca-certificates \
    && ln -s /opt/telegram-bot-isma/bin/telegram-bot /usr/local/bin/telegram-bot \
    && addgroup -g 1000 -S telegrambot && adduser -u 1000 -S telegrambot -G telegrambot 

USER telegrambot
