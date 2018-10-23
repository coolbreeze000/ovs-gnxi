FROM golang:alpine
ENV HOME=/home/target
RUN mkdir -p $HOME/certs
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR /ovs-gnxi
ADD ovs_gnxi /ovs-gnxi/ovs_gnxi
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/server.crt $HOME/certs/server.crt
ADD docker/target/certs/server.key $HOME/certs/server.key
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/ovs-gnxi/ovs_gnxi"]