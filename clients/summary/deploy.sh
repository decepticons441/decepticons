!/bin/bash
./build.sh

docker push nehay100/summary

export TLSCERT=/etc/letsencrypt/live/info441api.nehay.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/info441api.nehay.me/privkey.pem

# execute the upgrade-server.sh script on our server
# the -oStrictHostKeyChecking=no skips the prompt
# about adding the host to the list of known hosts
# so that this script doesn't get interrupted if the
# server's IP/hostname is new to us
ssh -oStrictHostKeyChecking=no root@$ip-address 'bash -s' < upgrade-server.sh