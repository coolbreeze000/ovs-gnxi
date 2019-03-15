#!/usr/bin/python

import socket
import argparse
import ipaddress

from mininet.net import Mininet
from mininet.cli import CLI
from mininet.log import setLogLevel
from mininet.node import Node, RemoteController, OVSSwitch


parser = argparse.ArgumentParser()
parser.add_argument("--controller", type=str, help="The hostname or IP address of the OpenFlow controller")
ARGS = parser.parse_args()


def resolve_controller(controller):
    try:
        ipaddress.ip_address(controller)
        return controller
    except ValueError:
        return socket.gethostbyname(controller)


def run_topology():
    # Initialize topology
    net = Mininet(controller=lambda name: RemoteController(
        name, ip=resolve_controller(ARGS.controller)), switch=OVSSwitch, autoSetMacs=True)
    # Add Controller
    net.addController( 'c0' )
    # Add hosts
    h1 = net.addHost('h1')
    h2 = net.addHost('h2')
    # Add switch
    sw1 = net.addSwitch('sw1', protocols='OpenFlow13', dpid='1')
    # Add links
    net.addLink(h1, sw1)
    net.addLink(h2, sw1)

    net.start()
    sw1.cmd('ovs-vsctl set-controller sw1 ssl:' + resolve_controller(ARGS.controller) + ':6653')
    CLI(net)
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    print('Deploying topology network...')
    run_topology()


topologies = {'ovs_gnxi_topology_network': (lambda: ToplogyNetwork())}
