#!/bin/bash

service openvswitch-switch stop
ovsdb-server /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info --remote=punix:/var/run/openvswitch/db.sock --private-key="/home/target/certs/target.key" --certificate="/home/target/certs/target.crt" --ca-cert "/home/target/certs/ca.crt" --no-chdir --log-file=/var/log/openvswitch/ovsdb-server.log --pidfile=/var/run/openvswitch/ovsdb-server.pid --detach --monitor
ovs-vswitchd unix:/var/run/openvswitch/db.sock -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log --pidfile=/var/run/openvswitch/ovs-vswitchd.pid --detach --monitor
ovs-appctl -t ovsdb-server ovsdb-server/add-remote pssl:6640
sleep 3