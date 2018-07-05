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

## Run

**Run Container:**
```bash
docker run -d -p 32123:32123 --name ovs-gnxi hub.docker.com/coolbreeze/ovs-gnxi/ovs-gnxi:latest
```

**Remove all stopped containers:**
```bash
docker ps -a -q | % { docker rm $_ }
```