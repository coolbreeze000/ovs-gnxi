FROM ubuntu:18.04
ENV HOME=/home/target
WORKDIR $HOME
RUN mkdir -p $HOME/certs
RUN mkdir -p /var/log/gnxi_target
RUN touch /var/log/gnxi_target/gnxi_target.log
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/target.crt $HOME/certs/target.crt
ADD docker/target/certs/target.key $HOME/certs/target.key
ADD docker/target/ovs_gnxi_topology_network.py /opt/
ADD docker/target/start_target.sh /home/target/start_target.sh
ADD docker/target/start_ovs.sh /home/target/start_ovs.sh
ADD docker/target/stop_ovs.sh /home/target/stop_ovs.sh
ADD docker/target/restart_ovs.sh /home/target/restart_ovs.sh
ADD target/gnxi_target $HOME/gnxi_target
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https net-tools iproute2 iputils-ping dnsutils openvswitch-common openvswitch-switch python-pip mininet screen git golang-go
RUN pip install ipaddress
CMD /home/target/start_target.sh && \
    /bin/bash
EXPOSE 10161/tcp
EXPOSE 10162/tcp
EXPOSE 18080/tcp
EXPOSE 32123/tcp