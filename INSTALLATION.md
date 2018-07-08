## Build

### Windows

```bash
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED=0
go build -o ovs_gnxi
```

### Mac OS X

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cod
```

## Push Container to Registry

```bash
docker login
docker pull coolbreeze/ovs-gnxi:ovs
docker pull coolbreeze/ovs-gnxi:gnxi
docker build -t coolbreeze/ovs-gnxi:ovs -f docker/ovs/Dockerfile .
docker build -t coolbreeze/ovs-gnxi:gnxi -f docker/gnxi/Dockerfile .
docker push coolbreeze/ovs-gnxi:ovs
docker push coolbreeze/ovs-gnxi:gnxi
```

## Docker-Compose

```bash
docker login
docker-compose pull
docker-compose build
docker-compose push
```

## Run

**Run Container:**
```bash
docker-compose up -d
docker-compose down
```

**Remove all stopped containers:**
```bash
docker ps -a -q | % { docker rm $_ }
```

**Remove all images:**
```bash
docker images -a -q | % { docker rmi $_ -f }
```

**Connect to Container:**
```bash
docker ps -aqf "name=gnxi"
docker exec -i -t $CONTAINER_ID bash
```

## Go

### Govendor

```bash
govendor init
govendor add +external
```
