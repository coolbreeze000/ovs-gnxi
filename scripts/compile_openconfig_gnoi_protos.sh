#!/bin/bash

proto_imports=".:${GOPATH}/src:${GOPATH}/src/ovs-gnxi/vendor:${GOPATH}/src/ovs-gnxi/shared/gnoi/modeldata/generated"

protoc -I=$proto_imports --go_out=plugins=grpc:./ ${GOPATH}/src/ovs-gnxi/shared/gnoi/modeldata/generated/system/system.proto
protoc -I=$proto_imports --go_out=plugins=grpc:./ ${GOPATH}/src/ovs-gnxi/shared/gnoi/modeldata/generated/cert/cert.proto
cp -f ovs-gnxi/shared/gnoi/modeldata/generated/system/system.pb.go ${GOPATH}/src/ovs-gnxi/shared/gnoi/modeldata/generated/system/system.pb.go
cp -f ovs-gnxi/shared/gnoi/modeldata/generated/cert/cert.pb.go ${GOPATH}/src/ovs-gnxi/shared/gnoi/modeldata/generated/cert/cert.pb.go
rm -rf ovs-gnxi