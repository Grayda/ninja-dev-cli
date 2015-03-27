FROM ubuntu:14.04
MAINTAINER Theo Julienne <theo@ninjablocks.com>

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:snappy-dev/beta
RUN apt-get update
RUN apt-get install -y snappy-tools bzr git build-essential ruby-dev
RUN gem install fpm

ADD scripts /scripts

ENTRYPOINT ["/scripts/entry"]
