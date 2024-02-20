FROM golang:1.21-alpine3.19

RUN go version
ENV GOPATH=/

# apk — это сокращение от Alpine Linux package manager (менеджер пакетов Alpine Linux)
RUN apk update && apk upgrade && apk add bash

# Нужно взять файлы и папки из локального контекста сборки и добавить их в текущую рабочую директорию образа
COPY ./ ./

# build go app
RUN go mod download
RUN go build -o myhttpserver ./cmd/main.go

CMD ["./myhttpserver"]