FROM golang:1.14-alpine AS build

# Обновление репозитория и установка git
RUN apk update && apk upgrade && apk add --no-cache git

# Переключение на /tmp/app в качестве рабочего каталога
WORKDIR /tmp/app

# Если имеются файлы go.mod или go.sum их нужно обязательно указывать
COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . .
# Сборка текущего проекта в бинарный файл с именем load_balancing_go
# Расположение бинарного файла по адресу /tmp/app/out/load_balancing_go
RUN GOOS=linux go build -o ./out/go_graphql_api .

FROM alpine:latest

# Добавляет CA Certificates в image
RUN apk add ca-certificates

# Копирует бинарный файл из BUILD контейнера в /app директорию
COPY --from=build /tmp/app/out/go_graphql_api /app/go_graphql_api

# Переключает в рабочую директорию /app
WORKDIR "/app"

# Выставляет порт 5000 из контейнера
EXPOSE 5000

# Запускает бинарник после запуска контейнера
CMD ["./go_graphql_api"]