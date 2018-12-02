FROM ubuntu:18.04
RUN mkdir -p /home/ovs/certs
ADD docker/ovs/certs/ca.crt /home/ovs/certs/ca.crt
ADD docker/ovs/certs/ovs.crt /home/ovs/certs/ovs.crt
ADD docker/ovs/certs/ovs.key /home/ovs/certs/ovs.key
ADD docker/ovs/ovs_gnxi_topology_network.py /opt/
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https net-tools iproute2 iputils-ping dnsutils openvswitch-common openvswitch-switch python-pip mininet screen
RUN pip install ipaddress
CMD service openvswitch-switch start --system-id=random && \
    service openvswitch-switch stop && \
    ovsdb-server /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info --remote=punix:/var/run/openvswitch/db.sock --private-key="/home/ovs/certs/ovs.key" --certificate="/home/ovs/certs/ovs.crt" --ca-cert "/home/ovs/certs/ca.crt" --no-chdir --log-file=/var/log/openvswitch/ovsdb-server.log --pidfile=/var/run/openvswitch/ovsdb-server.pid --detach --monitor && \
    ovs-vswitchd unix:/var/run/openvswitch/db.sock -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log --pidfile=/var/run/openvswitch/ovs-vswitchd.pid --detach --monitor && \
    ovs-appctl -t ovsdb-server ovsdb-server/add-remote pssl:6640 && \
    screen -S mininet -dm bash -c "python /opt/ovs_gnxi_topology_network.py --controller 'faucet.gnxi.lan' <<<'h1 ping h2'" && \
    /bin/bash