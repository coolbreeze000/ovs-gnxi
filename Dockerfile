FROM golang:alpine
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR /ovs-gnxi
#ADD ovsgnxi /ovs-gnxi/ovsgnxi
#EXPOSE 32123/tcp
#ENTRYPOINT ["/ovs-gnxi/ovsgnxi"]

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y apt-transport-https openvswitch-switch