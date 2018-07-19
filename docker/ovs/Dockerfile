FROM ubuntu:18.04
ADD docker/ovs/ovs_gnxi_topology_network.py /opt/
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https net-tools iputils-ping dnsutils openvswitch-common openvswitch-switch python-pip mininet
RUN pip install ipaddress
# RUN /usr/share/openvswitch/scripts/ovs-ctl start --system-id=random
# service openvswitch-switch start