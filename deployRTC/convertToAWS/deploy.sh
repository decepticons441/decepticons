#!/bin/sh

sh build.sh

docker push nestan/webrtc

ssh -i "webRTC.pem" ec2-user@ec2-3-209-165-5.compute-1.amazonaws.com 'bash-s' < upgrade-server.sh