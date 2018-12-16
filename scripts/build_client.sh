#!/bin/bash

echo "**Build Client**"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../client/gnxi_client ../client/