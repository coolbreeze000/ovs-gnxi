#!/bin/bash

service openvswitch-switch start --system-id=random
service openvswitch-switch stop
ovsdb-server /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info --remote=punix:/var/run/openvswitch/db.sock --private-key="/home/target/certs/active/target.key" --certificate="/home/target/certs/active/target.crt" --ca-cert "/home/target/certs/active/ca.crt" --no-chdir --log-file=/var/log/openvswitch/ovsdb-server.log --pidfile=/var/run/openvswitch/ovsdb-server.pid --detach --monitor
ovs-vswitchd unix:/var/run/openvswitch/db.sock -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log --pidfile=/var/run/openvswitch/ovs-vswitchd.pid --detach --monitor --private-key="/home/target/certs/active/target.key" --certificate="/home/target/certs/active/target.crt" --ca-cert "/home/target/certs/active/ca.crt"
ovs-appctl -t ovsdb-server ovsdb-server/add-remote pssl:6640
ovs-vsctl set-ssl /home/target/certs/active/target.key /home/target/certs/active/target.crt /home/target/certs/active/ca.crt
screen -S mininet -dm bash -c "python /opt/ovs_gnxi_topology_network.py --controller 'faucet.gnxi.lan' <<<'h1 ping h2'"
sleep 3
screen -S target -dm bash -c "/home/target/gnxi_target 2>&1 | tee /var/log/gnxi_target/gnxi_target.log"