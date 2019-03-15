FROM faucet/faucet:1.8.35
RUN mkdir -p /etc/faucet
RUN mkdir -p /etc/ryu/ssl
ADD docker/faucet/configs/ /etc/faucet/
ADD docker/faucet/certs/ca.crt  /etc/ryu/ssl/ca.crt
ADD docker/faucet/certs/faucet.crt /etc/ryu/ssl/faucet.crt
ADD docker/faucet/certs/faucet.key  /etc/ryu/ssl/faucet.key

ENTRYPOINT faucet --ryu-ctl-privkey "/etc/ryu/ssl/faucet.key" --ryu-ctl-cert "/etc/ryu/ssl/faucet.crt" --ryu-ca-certs "/etc/ryu/ssl/ca.crt"
