#!/bin/sh

echo "remove webrtc server"
docker rm -f webrtc

echo "pull docker image"
docker pull nestan/webrtc

echo "building webrtc server"
docker run -d \
--name webrtc \
-p 80:80 \
nestan/webrtc

echo "webrtc server built"