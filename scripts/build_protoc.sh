#!/bin/bash

cd protobuf-3.6.1
./configure
make
make check
make install
ldconfig
protoc --version