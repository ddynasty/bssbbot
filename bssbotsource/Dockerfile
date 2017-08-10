FROM golang:1.8.3

ADD . /go/src/github.com/bssbbot/bssbotsource

RUN go install github.com/bssbbot/bssbotsource

ENTRYPOINT /go/bin/bssbotsource

EXPOSE 8080