FROM ubuntu:14.04
MAINTAINER Theo Julienne <theo@ninjablocks.com>

ADD ubuntu-core-14.04-core-armhf.tar.gz /armhf

RUN apt-get update
RUN apt-get install -y qemu-user-static
RUN cp /usr/bin/qemu-arm-static /armhf/usr/bin

ADD go1.4.2.linux-amd64.tar.gz /usr/local
ADD go1.4.2.linux-arm~multiarch-armv7-1.tar.gz /armhf/usr/local

ADD scripts /scripts

RUN mkdir -p /armhf/go /go

RUN /scripts/go-patch-qemu

ENTRYPOINT ["/scripts/entry"]