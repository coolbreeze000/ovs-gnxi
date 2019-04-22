#!/bin/bash

# Source: https://github.com/faucetsdn/gnmi

OVSGNXI=$HOME/go/src/ovs-gnxi

# Generate CA Private Key
openssl req -newkey rsa:4096 -nodes -keyout ca.key -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=ca.gnxi.lan" 2> /dev/null
# Generate CA Certifiacate Signing Request
openssl req -key ca.key -new -out ca.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=ca.gnxi.lan" 2> /dev/null
# Generate and sign CA Certificate
openssl x509 -signkey ca.key -in ca.csr -req -days 365 -out ca.crt 2> /dev/null

# Generate Controller Private Key
openssl req -newkey rsa:4096 -nodes -keyout faucet.key -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=faucet.gnxi.lan" 2> /dev/null
# Generate Controller Certificate Signing Request
openssl req -key faucet.key -new -out faucet.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=faucet.gnxi.lan" 2> /dev/null
# Generate and sign Controller Certificate by CA
openssl x509 -req -in faucet.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out faucet.crt 2> /dev/null

# Generate Target Private Key
openssl req -newkey rsa:4096 -nodes -keyout target.key -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=target.gnxi.lan" 2> /dev/null
# Generate Target Certificate Signing Request
openssl req -key target.key -new -out target.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=target.gnxi.lan" 2> /dev/null
# Generate and sign Target Certificate by CA
openssl x509 -req -in target.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out target.crt 2> /dev/null

# Generate Client Private Key
openssl req -newkey rsa:4096 -nodes -keyout client.key -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=client.gnxi.lan" 2> /dev/null
# Generate Client Certificate Signing Request
openssl req -key client.key -new -out client.csr -subj "/C=AT/ST=Vienna/L=Vienna/O=Test/OU=Test/CN=client.gnxi.lan" 2> /dev/null
# Generate and sign Client Certificate by CA
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt 2> /dev/null

# Validate Controller Certificate
echo "**Validate Controller**"
openssl verify -verbose -CAfile ca.crt faucet.crt
# Validate Target Certificate
echo "**Validate Target**"
openssl verify -verbose -CAfile ca.crt target.crt
# Validate Client Certificate
echo "**Validate Client**"
openssl verify -verbose -CAfile ca.crt client.crt

# Remove unnecessary cert files
rm -f ca.csr faucet.csr target.csr client.csr
# Copy cert files to Controller
cp -f ./{ca.crt,faucet.crt,faucet.key} $OVSGNXI/docker/faucet/certs
# Copy cert files to Target
cp -f ./{ca.crt,target.crt,target.key} $OVSGNXI/docker/target/certs
# Copy cert files to Client
cp -f ./{ca.crt,ca.key,client.crt,client.key} $OVSGNXI/docker/client/certs
# Cleanup temporary cert files
rm -f ca.srl ca.crt ca.key faucet.crt faucet.key target.crt target.key client.crt client.key