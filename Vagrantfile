# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "ubuntu/bionic64"

  # Network Settings
  config.vm.network "private_network", type: "dhcp"

  # Machine Name
  config.vm.define "ovs-gnxi"

  # Config SSH Shell
  config.ssh.shell = "bash -c 'BASH_ENV=/etc/profile exec bash'"

  # Shared Folders
  config.vm.synced_folder "./", "/opt/ovs-gnxi/"

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL
    sudo -i
    hostname master.gnxi.lan
    echo "master.gnxi.lan" > /etc/hostname
    apt-get update -y
    apt-get install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    apt-get update -y && apt-get upgrade -y
    apt-get install -y docker-ce
    curl -sSL https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    apt-get install -y net-tools iputils-ping dnsutils git golang-go
    GOPATH=$HOME/go
    echo "GOPATH=$GOPATH" >> /etc/environment
    echo "GOBIN=$GOPATH/bin" >> /etc/environment
    echo "PATH=$GOBIN:${PATH}" >> /etc/environment
    mkdir -p "$GOPATH"
  SHELL
end
