#!/bin/bash

echo "**Build Target**"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../target/gnxi_target ../target/

echo "**Build Client**"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../client/gnxi_client ../client/