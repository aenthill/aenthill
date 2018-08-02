FROM golang:1.10-alpine3.7 AS builder

# Defines our working directory.
WORKDIR /go/src/github.com/aenthill/aenthill

# Copies our Go source.
COPY . . 

RUN go build -o aenthill cmd/aenthill/main.go

FROM alpine:3.7

# Installs missing libraries.
RUN apk add --no-cache --update \
    wget \
    tar
    
# Installs Docker client.
ENV DOCKER_VERSION "18.03.1-ce"
RUN wget -qO- https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz | tar xvz -C . &&\
    mv ./docker/docker /usr/bin &&\
    rm -rf ./docker

# Copies our aent entry point.
COPY build/test/docker/aent.sh /usr/bin/aent

# Copies aenthill binary.
COPY --from=builder /go/src/github.com/aenthill/aenthill/aenthill /usr/bin/aenthill