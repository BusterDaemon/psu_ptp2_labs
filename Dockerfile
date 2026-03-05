FROM go:1.26-alpine3.23-sfw-ent-dev AS build
WORKDIR /goida
COPY . .
#Запуск на сборку
#RUN command

#Контейнер alpine с собранным бинарником из прошлой стадии
#FROM alpine:3.23 AS main