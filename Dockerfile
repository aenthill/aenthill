FROM golang:1.10-stretch

LABEL authors="Julien Neuhart <j.neuhart@thecodingmachine.com>"

# Defines our working directory.
WORKDIR /go/src/github.com/anthill-docker/anthill

# Copies our Go source.
COPY . .

# Installs project dependencies.
RUN go get -d -v ./...

ENTRYPOINT [ "docker-entrypoint.sh" ]