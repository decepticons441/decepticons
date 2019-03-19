#!/bin/sh

sh build.sh

echo "pushing docker image"
docker push nestan/webrtc

echo "SSH into EC2"
ssh -i "webRTC.pem" ec2-user@ec2-3-209-165-5.compute-1.amazonaws.com 'bash -s' < upgrade-server.sh
