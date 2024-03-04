FROM golang:1.21

WORKDIR /go/api

COPY ./api/ /go/api
COPY ./common/ /go/common
COPY ./postgres/ /go/postgres
COPY ./redisdb/ /go/redisdb
COPY ./natsq/ /go/natsq

RUN go build main
EXPOSE 8080

ENTRYPOINT ["/go/api/main"]