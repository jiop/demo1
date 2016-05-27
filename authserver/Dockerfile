FROM golang:latest

MAINTAINER Jérémie IOP <jeremie.iop@gmail.com>

RUN mkdir /app

ADD src/main.go /app/main.go

WORKDIR /app

EXPOSE 8002

ENTRYPOINT ["go", "run", "main.go"]
