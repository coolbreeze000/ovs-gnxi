FROM ubuntu:18.04
ADD docker/ovs/ovs_gnxi_topology_network.py /opt/
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https net-tools iproute2 iputils-ping dnsutils openvswitch-common openvswitch-switch python-pip mininet screen
RUN pip install ipaddress
CMD service openvswitch-switch start --system-id=random && \
    ovs-appctl -t ovsdb-server ovsdb-server/add-remote ptcp:6640 && \
    screen -S mininet -dm bash -c "python /opt/ovs_gnxi_topology_network.py --controller 'faucet.gnxi.lan'" && \
    /bin/bash