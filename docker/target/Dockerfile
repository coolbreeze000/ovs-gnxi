FROM golang:alpine
ENV HOME=/home/target
RUN mkdir -p $HOME/certs
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR $HOME
ADD target/gnxi_target $HOME/gnxi_target
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/target.crt $HOME/certs/target.crt
ADD docker/target/certs/target.key $HOME/certs/target.key
ADD docker/target/openconfig-openflow.json $HOME/openconfig-openflow.json
EXPOSE 10161/tcp
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/home/target/gnxi_target"]