FROM ubuntu:18.04
ENV HOME=/home/client
RUN mkdir -p $HOME/certs
WORKDIR $HOME
ADD client/gnxi_client $HOME/gnxi_client
ADD docker/client/certs/ca.crt $HOME/certs/ca.crt
ADD docker/client/certs/ca.key $HOME/certs/ca.key
ADD docker/client/certs/client.crt $HOME/certs/client.crt
ADD docker/client/certs/client.key $HOME/certs/client.key
ADD docker/target/certs/target.crt $HOME/certs/target.crt
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -qy --no-install-recommends apt-utils ca-certificates net-tools iputils-ping dnsutils openvswitch-common git golang-go
ENV GOPATH=$HOME/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:${PATH}
RUN mkdir -p $GOPATH
ENV GNMI_TARGET=target.gnxi.lan
ENV GNMI_PORT=32123
