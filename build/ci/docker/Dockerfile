FROM golang:1.10-stretch

# Installs Docker client.
ENV DOCKER_VERSION "18.03.1-ce"
RUN wget -qO- https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz | tar xvz -C . &&\
    mv ./docker/docker /usr/bin &&\
    rm -rf ./docker

# Installs lint dependencies.
RUN go get -u gopkg.in/alecthomas/gometalinter.v2 &&\
    gometalinter.v2 --install

# Defines our working directory.
WORKDIR /go/src/github.com/aenthill/aenthill

# Copies our Go source.
COPY . . 

ENTRYPOINT [ "./build/ci/docker/docker-entrypoint.sh" ]