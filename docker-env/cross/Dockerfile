FROM ubuntu:14.04
MAINTAINER Theo Julienne <theo@ninjablocks.com>

RUN apt-get update
RUN apt-get install -y click click-dev schroot dpkg-dev debootstrap mercurial git golang-go build-essential

ADD scripts-cross /scripts-cross

ENTRYPOINT ["/scripts-cross/entry"]
