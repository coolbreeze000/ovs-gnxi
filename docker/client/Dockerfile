# Source: https://github.com/faucetsdn/gnmi

FROM ubuntu:18.04
ENV HOME=/home/client
RUN mkdir $HOME
WORKDIR $HOME
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -qy --no-install-recommends apt-utils ca-certificates net-tools iputils-ping dnsutils git golang-go
ENV GOPATH=$HOME/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:${PATH}
RUN mkdir -p $HOME/certs
RUN mkdir -p $GOPATH
RUN go get -u github.com/google/gnxi/gnmi_capabilities
RUN go get -u github.com/google/gnxi/gnmi_get
RUN go get -u github.com/google/gnxi/gnmi_set
RUN go install -v github.com/google/gnxi/gnmi_capabilities
RUN go install -v github.com/google/gnxi/gnmi_get
RUN go install -v github.com/google/gnxi/gnmi_set

RUN cd $HOME/certs
# Generate CA Private Key
RUN openssl req -newkey rsa:4096 -nodes -keyout ca.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate CA Certifiacate Signing Request
RUN openssl req -key ca.key -new -out ca.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate and sign CA Certificate
RUN openssl x509 -signkey ca.key -in ca.csr -req -days 365 -out ca.crt
# Generate Client Private Key
RUN openssl req -newkey rsa:4096 -nodes -keyout client.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate Client Certificate Signing Request
RUN openssl req -key client.key -new -out client.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate and sign Client Certificate by CA
RUN openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt
# Validate Client Certificate
RUN echo "**Validate Client**"
RUN openssl verify -verbose -CAfile ca.crt client.crt

ENV GNMI_TARGET=target.gnxi.lan
ENV GNMI_PORT=32123