# Стадия сборки
FROM golang:latest AS builder

WORKDIR /Habr-comments-server
COPY . .

RUN go mod tidy && go build -o server ./cmd/server/main.go

# Стадия запуска
FROM ubuntu:22.04

RUN apt update && apt install -y make

WORKDIR /Habr-comments-server
COPY --from=builder /Habr-comments-server/server /Habr-comments-server/server
COPY ./config ./config
COPY ./dev.env ./dev.env

ENV CONFIG_PATH=./config/dev.yaml

EXPOSE 8082

CMD ["./server"]
