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
export OVS_CONTAINER_ID=`docker ps -aqf 'name=client'`
docker exec -i -t $OVS_CONTAINER_ID bash
```

#### Container: OVS (mininet)
```bash
screen -ls
screen -r mininet
# CTRL+a d
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
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ovs_gnxi
docker-compose up -d --force-recreate --build target
docker logs target
```

### Test Client
```bash
gnmi_get -target_addr target:10161 -key certs/client.key -cert certs/client.crt -ca certs/ca.crt -target_name server.gnxi.lan -alsologtostderr \
  -xpath "/system/openflow/agent/config/datapath-id" \
  -xpath "/system/openflow/controllers/controller[name=main]/connections/connection[aux-id=0]/config/address"
```