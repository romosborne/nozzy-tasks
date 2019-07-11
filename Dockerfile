FROM golang:1.11

WORKDIR $GOPATH/src/github.com/romosborne/nozzy-tasks
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

VOLUME /config

CMD ["nozzy-tasks"]
