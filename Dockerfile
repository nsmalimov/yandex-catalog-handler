FROM golang:latest

RUN mkdir /app

ADD . /app

WORKDIR /app/cmd

RUN go build -o main .

EXPOSE 8080

CMD ["./main", "-config-path", "/app/configs/config.yaml"]