# !/bin/bash
export PRIVNETNAME="privnet"

docker rm -f summary

# ADDING MAIN.GO CODE
docker pull nehay100/summary

docker run -d \
--network $PRIVNETNAME \
-e ADDR="summary:80" \
--name summary \
nehay100/summary

