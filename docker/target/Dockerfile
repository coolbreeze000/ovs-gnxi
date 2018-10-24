FROM golang:alpine
ENV HOME=/home/target
RUN mkdir -p $HOME/certs
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR $HOME
ADD ovs_gnxi $HOME/ovs_gnxi
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/server.crt $HOME/certs/server.crt
ADD docker/target/certs/server.key $HOME/certs/server.key
ADD docker/target/openconfig-openflow.json $HOME/openconfig-openflow.json
EXPOSE 10161/tcp
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/home/target/ovs_gnxi"]