FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN apt-get update && apt-get install -y default-mysql-client
RUN go get -d -v ./...
RUN go install -v ./...

CMD database-backup-manager