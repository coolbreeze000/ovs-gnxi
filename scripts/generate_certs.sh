#!/bin/bash

# Source: https://github.com/faucetsdn/gnmi

OVSGNXI=$HOME/go/src/ovs-gnxi

# Generate CA Private Key
openssl req -newkey rsa:4096 -nodes -keyout ca.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate CA Certifiacate Signing Request
openssl req -key ca.key -new -out ca.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate and sign CA Certificate
openssl x509 -signkey ca.key -in ca.csr -req -days 365 -out ca.crt

# Generate Target Private Key
openssl req -newkey rsa:4096 -nodes -keyout target.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=target.gnxi.lan"
# Generate Target Certificate Signing Request
openssl req -key target.key -new -out target.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=target.gnxi.lan"
# Generate and sign Target Certificate by CA
openssl x509 -req -in target.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out target.crt

# Generate Client Private Key
openssl req -newkey rsa:4096 -nodes -keyout client.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate Client Certificate Signing Request
openssl req -key client.key -new -out client.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate and sign Client Certificate by CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt

# Validate Target Certificate
echo "**Validate Target**"
openssl verify -verbose -CAfile ca.crt target.crt
# Validate Client Certificate
echo "**Validate Client**"
openssl verify -verbose -CAfile ca.crt client.crt

# Remove unnecessary cert files
rm -f ca.key ca.csr target.csr client.csr
# Copy cert files to Target
cp -f ./{ca.crt,target.crt,target.key} $OVSGNXI/docker/target/certs
# Copy cert files to Client
cp -f ./{ca.crt,client.crt,client.key} $OVSGNXI/docker/client/certs
# Cleanup temporary cert files
rm -f ca.srl ca.crt target.crt target.key client.crt client.key