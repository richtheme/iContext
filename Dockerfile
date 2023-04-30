FROM golang:1.20.3-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# build go app
RUN go mod download
RUN go build -o iContext ./cmd/api/main.go

CMD ["./iContext"]
EXPOSE 8080/tcp 