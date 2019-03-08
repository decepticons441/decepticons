docker rm -f nehay100/summary
docker pull nehay100/summary

docker run -d \
-p 443:443 \
-p 80:80 \
-e TLSKEY=$TLSKEY
-e TLSCERT=$TLSCERT
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--name summary
nehay100/summary

exit
EOF