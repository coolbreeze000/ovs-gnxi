FROM golang:alpine
RUN apk add --no-cache tzdata
ENV TZ=Europe/Vienna
WORKDIR /ovs-gnxi
ADD ovs_gnxi /ovs-gnxi/ovs_gnxi
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/ovs-gnxi/ovs_gnxi"]