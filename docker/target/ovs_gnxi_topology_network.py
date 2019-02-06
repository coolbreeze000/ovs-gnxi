#!/usr/bin/python

import socket
import argparse
import ipaddress

from mininet.topo import Topo
from mininet.net import Mininet
from mininet.cli import CLI
from mininet.log import setLogLevel
from mininet.node import Node, RemoteController, OVSSwitch


parser = argparse.ArgumentParser()
parser.add_argument("--controller", type=str, help="The hostname or IP address of the OpenFlow controller")
ARGS = parser.parse_args()


class ToplogyNetwork(Topo):

    def __init__(self):
        # Initialize topology
        Topo.__init__(self)

        # Add hosts
        h1 = self.addHost('h1')
        h2 = self.addHost('h2')

        # Add switches
        sw1 = self.addSwitch('sw1', protocols='OpenFlow13', dpid='1')

        # Add links
        self.addLink(h1, sw1)
        self.addLink(h2, sw1)


def resolve_controller(controller):
    try:
        ipaddress.ip_address(controller)
        return controller
    except ValueError:
        return socket.gethostbyname(controller)


def run_topology():
    topology = ToplogyNetwork()
    net = Mininet(topo=topology, controller=lambda name: RemoteController(
        name, ip=resolve_controller(ARGS.controller)), switch=OVSSwitch, autoSetMacs=True)
    net.start()
    CLI(net)
    net.stop()


if __name__ == '__main__':
    setLogLevel('info')
    print('Deploying topology network...')
    run_topology()


topologies = {'ovs_gnxi_topology_network': (lambda: ToplogyNetwork())}
