#!/bin/sh
docker rm -f webrtc

docker pull nestan/webrtc

echo "removed webrtc server"

docker run -d \
--name webrtc \
-p 80:80 \
nestan/webrtc