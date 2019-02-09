#!/bin/bash

cd /root/protobuf-3.6.1
./configure
make
make install
ldconfig
protoc --version