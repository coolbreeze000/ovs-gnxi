#!/bin/bash

# Source: https://github.com/faucetsdn/gnmi

# Generate CA Private Key
openssl req -newkey rsa:4096 -nodes -keyout ca.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate CA Certifiacate Signing Request
openssl req -key ca.key -new -out ca.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=ca.gnxi.lan"
# Generate and sign CA Certificate
openssl x509 -signkey ca.key -in ca.csr -req -days 365 -out ca.crt

# Generate Server Private Key
openssl req -newkey rsa:4096 -nodes -keyout server.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=server.gnxi.lan"
# Generate Server Certificate Signing Request
openssl req -key server.key -new -out server.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=server.gnxi.lan"
# Generate and sign Server Certificate by CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

# Generate Client Private Key
openssl req -newkey rsa:4096 -nodes -keyout client.key -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate Client Certificate Signing Request
openssl req -key client.key -new -out client.csr -subj "/C=AT/ST=Vienna/L=Test/O=Test/OU=Test/CN=client.gnxi.lan"
# Generate and sign Client Certificate by CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt

# Validate Server Certificate
echo "**Validate Server**"
openssl verify -verbose -CAfile ca.crt server.crt
# Validate Client Certificate
echo "**Validate Client**"
openssl verify -verbose -CAfile ca.crt client.crt

# Remove unnecessary cert files
rm -f ca.key ca.csr server.csr client.csr
# Copy cert files to target
cp ./{ca.crt,server.crt,server.key} docker/target/certs
# Copy cert files to client
cp ./{ca.crt,client.crt,client.key} docker/client/certs
# Cleanup temporary cert files
rm -f ca.crt server.crt server.key ca.crt client.crt client.key