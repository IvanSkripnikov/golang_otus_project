FROM golang:1.19-alpine

WORKDIR /go/src/app

RUN apk update && \
    apk add git && \
    apk add --no-cache alpine-conf

ADD ./app /go/src/app

COPY ./.env.example ./.env

RUN cd /go/src/app && \
    go install github.com/IvanSkripnikov/golang_otus_project

WORKDIR "/go"
