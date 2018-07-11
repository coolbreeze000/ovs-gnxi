## Build

### On Windows

#### For Windows

```bash
$env:GOOS="windows"
$env:GOARCH="amd64"
$env:CGO_ENABLED=0
set PATHEXT=.BAT
go build -o ovs_gnxi
```

#### For Linux

```bash
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED=0
set PATHEXT=.BAT
go build -o ovs_gnxi
```

### On Mac OS X

#### For Linux

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cod
```

### Windows

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o cod
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

### Run Container
```bash
docker-compose build
docker-compose up -d --force-recreate
docker-compose run ovs python /opt/ovs_gnxi_topology_network.py --controller 'faucet.gnxi.lan'
docker-compose down
```

### Stop all containers
```bash
docker ps -a -q | % { docker stop $_ }
```

### Remove all stopped containers
```bash
docker ps -a -q | % { docker rm $_ }
```

### Remove all images
```bash
docker images -a -q | % { docker rmi $_ -f }
```

### Connect to Container

#### Windows

```bash
$env:OVS_CONTAINER_ID=docker ps -aqf "name=ovs"
docker exec -i -t $env:OVS_CONTAINER_ID bash
```

## Go

### Govendor

```bash
govendor init
govendor add +external
```
