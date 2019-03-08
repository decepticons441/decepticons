# !/bin/bash
sh build.sh

docker push nehay100/gateway

# ssh root@ip-address

# execute the upgrade-server.sh script on our server
# the -oStrictHostKeyChecking=no skips the prompt
# about adding the host to the list of known hosts
# so that this script doesn't get interrupted if the
# server's IP/hostname is new to us
ssh -oStrictHostKeyChecking=no ec2-user@3.17.90.152 'bash -s' < upgrade-server.sh