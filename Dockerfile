# Base usada para o container "builder"
FROM golang:1.15-alpine as builder

# Env var que faz o go usar go.mod
ENV GO111MODULE on

# Dependencias básicas
RUN apk --no-cache add git make tree

# Diretório padrão para trabalho
WORKDIR /home/src/github.com/ismaelpereira/bot-telegram-isma/

# Copia só os arquivos de dependencias primeiro (otimiza cache do build)
COPY go.mod .
COPY go.sum .

# Roda o go mod download para baixar dependencias (otimiza cache do build)
RUN go mod download

# Copia o resto dos arquivos
COPY . .

# faz o build
RUN make build

###

# Base usada para a imagem final
FROM alpine:3.13

# Copia os arquivos da imagem "builder" para a imagem final
COPY --from=builder /home/src/github.com/ismaelpereira/bot-telegram-isma/dist/config/ /etc/telegram-bot/
COPY --from=builder /home/src/github.com/ismaelpereira/bot-telegram-isma/dist/bin/ /opt/telegram-bot/bin/

# Diretório padrão para trabalho
WORKDIR /opt/telegram-bot/

# instala ca-certificates (faz https (ssl) funcionar)
# cria usuario e grupo de ID 1000 (evitar usar root pra container é uma boa ideia)
RUN apk --no-cache add ca-certificates \
      && addgroup -g 1000 -S telegrambot && adduser -u 1000 -S telegrambot -G telegrambot 

# definie o usuário padrão de execução dos comandos no container
USER telegrambot

# comando padrão do container
CMD ["/opt/telegram-bot/bin/telegram-bot","/etc/telegram-bot/default.json"]
