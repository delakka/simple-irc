FROM golang:1.12-alpine
ADD . /go/src/simple-irc
WORKDIR /go/src/simple-irc
EXPOSE 6667
CMD [ "go", "run", "simple-irc" ]