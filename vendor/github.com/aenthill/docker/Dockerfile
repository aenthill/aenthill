FROM golang:1.10-stretch

LABEL authors="Julien Neuhart <j.neuhart@thecodingmachine.com>"

# Installs lint dependencies.
RUN go get -u gopkg.in/alecthomas/gometalinter.v2 &&\
    gometalinter.v2 --install

# Defines our working directory.
WORKDIR /go/src/github.com/aenthill/docker

# Copies our Go source.
COPY . .

# Installs project dependencies.
RUN go get -d -v ./...

# Installs Docker client.
ENV DOCKER_VERSION "18.03.1-ce"
RUN wget -qO- https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz | tar xvz -C . &&\
    mv ./docker/docker /usr/bin &&\
    rm -rf ./docker

# Defines SHELL.
ENV SHELL "/bin/sh"

ENTRYPOINT [ "./docker-entrypoint.sh" ]