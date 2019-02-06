FROM ubuntu:18.04
ENV HOME=/home/target
WORKDIR $HOME
RUN mkdir -p $HOME/certs
ADD docker/target/certs/ca.crt $HOME/certs/ca.crt
ADD docker/target/certs/target.crt $HOME/certs/target.crt
ADD docker/target/certs/target.key $HOME/certs/target.key
ADD docker/target/ovs_gnxi_topology_network.py /opt/
ADD target/gnxi_target $HOME/gnxi_target
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y apt-transport-https net-tools iproute2 iputils-ping dnsutils openvswitch-common openvswitch-switch python-pip mininet screen git golang-go
RUN pip install ipaddress
RUN service openvswitch-switch start --system-id=random && \
    service openvswitch-switch stop && \
    ovsdb-server /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info --remote=punix:/var/run/openvswitch/db.sock --private-key="/home/target/certs/target.key" --certificate="/home/target/certs/target.crt" --ca-cert "/home/target/certs/ca.crt" --no-chdir --log-file=/var/log/openvswitch/ovsdb-server.log --pidfile=/var/run/openvswitch/ovsdb-server.pid --detach --monitor && \
    ovs-vswitchd unix:/var/run/openvswitch/db.sock -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log --pidfile=/var/run/openvswitch/ovs-vswitchd.pid --detach --monitor && \
    ovs-appctl -t ovsdb-server ovsdb-server/add-remote pssl:6640 && \
    screen -S mininet -dm bash -c "python /opt/ovs_gnxi_topology_network.py --controller 'faucet.gnxi.lan' <<<'h1 ping h2'"
EXPOSE 10161/tcp
EXPOSE 18080/tcp
EXPOSE 32123/tcp
ENTRYPOINT ["/home/target/gnxi_target"]