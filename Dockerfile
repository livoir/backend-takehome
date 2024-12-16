FROM golang:1.21.0

ENV GIN_MODE release

WORKDIR /go/src/app

RUN go install github.com/air-verse/air@v1.61.1

COPY ./app .

CMD air