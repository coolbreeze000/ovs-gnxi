## Build

### For Linux

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ovs_gnxi
```

## Login to Registry

```bash
docker login
```

## Pull Containers from Registry

```bash
docker pull coolbreeze/ovs-gnxi:ovs
docker pull coolbreeze/ovs-gnxi:target
docker pull coolbreeze/ovs-gnxi:client
```

## Build Containers

```bash
docker build -t coolbreeze/ovs-gnxi:ovs -f docker/ovs/Dockerfile .
docker build -t coolbreeze/ovs-gnxi:target -f docker/target/Dockerfile .
docker build -t coolbreeze/ovs-gnxi:client -f docker/client/Dockerfile .
```

## Push Containers to Registry

```bash
docker push coolbreeze/ovs-gnxi:ovs
docker push coolbreeze/ovs-gnxi:target
docker push coolbreeze/ovs-gnxi:client
```

## Docker-Compose

```bash
docker login
docker-compose pull
docker-compose build
docker-compose push
```

## Run

### Run Vagrant
```powershell
$env:VAGRANT_DEFAULT_PROVIDER="virtualbox"
vagrant box update
vagrant destroy --force
vagrant up
vagrant global-status
vagrant reload --provision
vagrant ssh ovs-gnxi
```


### Run Container
```bash
docker-compose build --no-cache
docker-compose up -d --force-recreate --build
docker-compose down
```

### Stop all containers
```bash
docker stop $(docker ps -a -q)
```

### Remove all stopped containers
```bash
docker rm $(docker ps -a -q)
```

### Remove all images
```bash
docker rmi $(docker images -a -q) -f
```

### Connect to Container

```bash
export CONTAINER_ID_OVS=`docker ps -aqf 'name=ovs'`
export CONTAINER_ID_FAUCET=`docker ps -aqf 'name=faucet'`
export CONTAINER_ID_TARGET=`docker ps -aqf 'name=target'`
export CONTAINER_ID_CLIENT=`docker ps -aqf 'name=client'`
docker exec -i -t $CONTAINER_ID_OVS bash
```

#### Container: OVS (mininet)
```bash
screen -ls
screen -r mininet
# CTRL+a d
ovs-ofctl -O OpenFlow13 dump-flows sw1
```

#### Container: Client (ovs-client)
```bash
ovsdb-client list-dbs tcp:`dig ovs.gnxi.lan +short`:6640
ovsdb-client list-tables
ovsdb-client list-columns Open_vSwitch Controller
ovsdb-client dump Open_vSwitch Controller
ovsdb-client dump Open_vSwitch Controller target
ovsdb-client transact '["Open_vSwitch",{"op":"select", "table":"Controller", "where": [], "columns":["target"]}]'
ovsdb-client transact '["Open_vSwitch",{"op":"select", "table":"Interface", "where": [], "columns":["name", "mtu", "link_state", "statistics"]}]'
ovsdb-client transact '["Open_vSwitch",{"op":"select", "table":"Open_vSwitch", "where": [], "columns":["ovs_version", "external_ids"]}]'
ovsdb-client transact '["Open_vSwitch",{"op":"update", "table":"Controller", "where":[["_uuid", "==", ["uuid","ebe661a9-7c45-4430-9b02-01992d96a1f0"]]], "row":{"target":"tcp:172.18.0.4:7777"}}]'
ovsdb-client transact '["Open_vSwitch",{"op":"update", "table":"Open_vSwitch", "where":[["_uuid", "==", ["uuid","017dfc26-c939-417f-99ce-e431b9c0d707"]]], "row":{"db_version":"8.8.8"}}]'
```

#### Container: OVS (startup)
```bash
ovsdb-tool create
ovsdb-server /etc/openvswitch/conf.db -vconsole:emer -vsyslog:err -vfile:info --system-id=random --remote=punix:/var/run/openvswitch/db.sock --private-key="/home/ovs/certs/ovs.key" --certificate="/home/ovs/certs/ovs.crt" --ca-cert "/home/ovs/certs/ca.crt" --no-chdir --log-file=/var/log/openvswitch/ovsdb-server.log --pidfile=/var/run/openvswitch/ovsdb-server.pid --detach --monitor
ovs-vswitchd unix:/var/run/openvswitch/db.sock -vconsole:emer -vsyslog:err -vfile:info --mlockall --no-chdir --log-file=/var/log/openvswitch/ovs-vswitchd.log --pidfile=/var/run/openvswitch/ovs-vswitchd.pid --detach --monitor
ovs-appctl -t ovsdb-server ovsdb-server/add-remote pssl:6640
ovsdb-client --private-key "/home/client/certs/client.key" --certificate "/home/client/certs/client.crt" --ca-cert "/home/client/certs/ca.crt" list-dbs ssl:`dig ovs.gnxi.lan +short`:6640
```

## Go

### Install Govendor

```powershell
go get -u github.com/kardianos/govendor
```

### Use Govendor

```powershell
govendor init
govendor add +external
```

## Development

### Sync Changes for Development

```bash
vagrant rsync-auto
vagrant ssh ovs-gnxi
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gnxi_target
docker-compose up -d --force-recreate --build target
docker logs target
```

### Test Client
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gnxi_client
docker-compose up -d --force-recreate --build client
gnmi_get -target_addr target:10161 -key certs/client.key -cert certs/client.crt -ca certs/ca.crt -target_name server.gnxi.lan -alsologtostderr \
  -xpath "/system/openflow/agent/config/datapath-id" \
  -xpath "/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address"
```

### Generate Go Bindings for Open vSwitch
go run vendor/github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=oc.go -package_name=openvswitch -generate_fakeroot -fakeroot_name=device -compress_paths=true -exclude_modules=ietf-interfaces yang/openconfig-interfaces.yang