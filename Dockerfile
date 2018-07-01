FROM golang:1.10.3 as builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . /go/src/github.com/mlabouardy/mon-put-instance-data
WORKDIR /go/src/github.com/mlabouardy/mon-put-instance-data

RUN dep ensure
RUN go build

FROM alpine

COPY --from=builder /go/src/github.com/mlabouardy/mon-put-instance-data/mon-put-instance-data /usr/bin/mon-put-instance-data

RUN apk update && apk add ca-certificates curl

RUN chmod +x /usr/bin/mon-put-instance-data && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

ENV DOCKER_VERSION "18.03.1-ce"
RUN curl -L -o /tmp/docker-$DOCKER_VERSION.tgz https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_VERSION.tgz && \
    tar -xz -C /tmp -f /tmp/docker-$DOCKER_VERSION.tgz && \
    mv /tmp/docker/docker /usr/bin && \
    rm -rf /tmp/docker-$VERSION /tmp/docker && \
    chmod +x /usr/bin/docker

ENTRYPOINT ["mon-put-instance-data"]