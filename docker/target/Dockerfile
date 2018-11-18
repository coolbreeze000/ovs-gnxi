FROM golang:alpine
ENV HOME=/home/target
RUN mkdir -p $HOME/certs
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR $HOME
ADD target/gnxi_target $HOME/gnxi_target
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/server.crt $HOME/certs/server.crt
ADD docker/target/certs/server.key $HOME/certs/server.key
ADD docker/target/openconfig-openflow.json $HOME/openconfig-openflow.json
EXPOSE 10161/tcp
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/home/target/gnxi_target"]