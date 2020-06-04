FROM golang:latest

COPY . .

RUN apt-get update
RUN apt-get install vim -y
RUN go get "github.com/go-sql-driver/mysql"