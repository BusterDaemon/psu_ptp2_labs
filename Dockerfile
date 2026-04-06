FROM golang:alpine3.23 AS build
WORKDIR /goida
COPY . .
RUN apk add gcc musl-dev
#Запуск на сборку
RUN go build -o ./apishka ./cmd/api/main.go

#Контейнер alpine с собранным бинарником из прошлой стадии
FROM alpine:3.23 AS main
WORKDIR /app
COPY --from=build /goida/apishka ./app
CMD [ "./app" ]