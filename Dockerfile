FROM golang:1.10-stretch

LABEL authors="Julien Neuhart <j.neuhart@thecodingmachine.com>"

# Installs lint dependencies.
RUN go get -u gopkg.in/alecthomas/gometalinter.v2 &&\
    gometalinter.v2 --install

# Defines our working directory.
WORKDIR /go/src/github.com/anthill-docker/anthill

# Copies our Go source.
COPY . .

# Installs project dependencies.
RUN go get -d -v ./...

ENTRYPOINT [ "./docker-entrypoint.sh" ]