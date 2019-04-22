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
  config.vm.box_version = "20190205.0.0"

  config.disksize.size = '20GB'

  # Network Settings
  config.vm.network "private_network", type: "dhcp"

  # Machine Name
  config.vm.define "ovs-gnxi"

  # Config SSH Shell
  config.ssh.shell = "bash -c 'BASH_ENV=/etc/profile exec bash'"

  # Shared Folders
  config.vm.synced_folder "./", "/root/go/src/ovs-gnxi/",
    type: "rsync",
    rsync__auto: true,
    rsync__exclude: [".idea/", ".vagrant/", ".github/", ".git/", ".travis.yml", "Vagrantfile", ".gitignore", "*.log", "*.crt", "*.csr", "*.key", "gnxi_target", "gnxi_client"]

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
    apt-get install -y autoconf automake libtool curl make g++ unzip protobuf-compiler golang-goprotobuf-dev
    cd /root
    wget -q https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-all-3.6.1.tar.gz
    tar xzf protobuf-all-3.6.1.tar.gz
    cd /root/go/src/ovs-gnxi/scripts/
    ./generate_certs.sh
    ./build_all.sh > /dev/null 2>&1 &
    cd ..
    docker-compose up -d --force-recreate --build
    cd
  SHELL
end
