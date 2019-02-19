# ovs-gnxi

A simple gNXI wrapper prototype for Open vSwitch (OVS).

## Run

```bash
vagrant destroy --force
vagrant up
vagrant rsync-auto
vagrant ssh ovs-gnxi
cd go/src/ovs-gnxi/scripts

./build_client.sh && docker-compose up -d --force-recreate --build client && export CONTAINER_ID_CLIENT=`docker ps -aqf 'name=client'` && docker exec -i -t $CONTAINER_ID_CLIENT bash
./gnxi_client

./build_target.sh && docker-compose up -d --force-recreate --build target && export CONTAINER_ID_TARGET=`docker ps -aqf 'name=target'` && docker exec -i -t $CONTAINER_ID_TARGET bash
tail -f /var/log/gnxi_target/gnxi_target.log
```

## Generate Go Bindings for Open vSwitch
go run vendor/github.com/openconfig/ygot/generator/generator.go -path=yang -output_file=oc.go -package_name=openvswitch -generate_fakeroot -fakeroot_name=device -compress_paths=true -exclude_modules=ietf-interfaces yang/openconfig-interfaces.yang