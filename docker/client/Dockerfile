FROM ubuntu:18.04
ENV HOME=/home/client
RUN mkdir $HOME/certs
WORKDIR $HOME
ADD docker/client/certs/ca.crt $HOME/certs/ca.crt
ADD docker/client/certs/client.crt $HOME/certs/client.crt
ADD docker/client/certs/client.key $HOME/certs/client.key
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -qy --no-install-recommends apt-utils ca-certificates net-tools iputils-ping dnsutils git golang-go
ENV GOPATH=$HOME/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:${PATH}
RUN mkdir -p $GOPATH
RUN go get -u github.com/google/gnxi/gnmi_capabilities
RUN go get -u github.com/google/gnxi/gnmi_get
RUN go get -u github.com/google/gnxi/gnmi_set
RUN go install -v github.com/google/gnxi/gnmi_capabilities
RUN go install -v github.com/google/gnxi/gnmi_get
RUN go install -v github.com/google/gnxi/gnmi_set
ENV GNMI_TARGET=target.gnxi.lan
ENV GNMI_PORT=32123
