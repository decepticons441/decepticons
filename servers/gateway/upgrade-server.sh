# !/bin/bash
export TLSCERT=/etc/letsencrypt/live/api.nehay.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.nehay.me/privkey.pem
# generate 18 random bytes and base64 encode them
export MYSQL_ROOT_PASSWORD="password"
# export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
# export SESSIONKEY=$(openssl rand -base64 18)
export SESSIONKEY="sessionkey"
export MYSQL_CONTAINER_NAME="mysqldemo"
export MYSQL_DB="userDB"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp($MYSQL_CONTAINER_NAME:3306)/$MYSQL_DB"
export REDISCLIENT="redisserver"
export REDISADDR="$REDISCLIENT:6379"
export PRIVNETNAME="privnet"
export MESSAGE_CONTAINER="node"
export SUMMARY_CONTAINER="summary"
export RABBIT="rabbit"

# add message addr and summary addr (being more than one addr)
export MESSAGE_ADDR="http://$MESSAGE_CONTAINER:80"
export SUMMARY_ADDR="http://$SUMMARY_CONTAINER:80"

# docker network disconnect -f $PRIVNETNAME $MYSQL_CONTAINER_NAME
# docker network disconnect -f $PRIVNETNAME $REDISCLIENT
# docker network disconnect -f $PRIVNETNAME gateway
docker rm -f $MYSQL_CONTAINER_NAME
docker rm -f $REDISCLIENT
docker rm -f $MESSAGE_CONTAINER
docker rm -f $SUMMARY_CONTAINER
docker rm -f rabbit
docker rm -f gateway

docker network rm $PRIVNETNAME

# CREATE PRIVATE NETWORK
docker network create $PRIVNETNAME
 
docker image prune -f
docker container prune -f
docker volume prune -f

docker pull rabbitmq:3

docker run -d \
--network $PRIVNETNAME \
--name rabbit \
rabbitmq:3

# ADDING MYSQL CONTAINER TO PRIVATE NETWORK HOST
docker pull nehay100/$MYSQL_CONTAINER_NAME

docker run -d \
--network $PRIVNETNAME \
--name $MYSQL_CONTAINER_NAME \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=userDB \
nehay100/$MYSQL_CONTAINER_NAME

sleep 20

# ADDING REDIS TO PRIVATE NETWORK HOST
docker run -d \
--network $PRIVNETNAME \
--name $REDISCLIENT \
redis

# ADDING MAIN.GO CODE
docker pull nehay100/gateway

docker run -d \
--network $PRIVNETNAME \
-p 443:443 \
-e TLSKEY=$TLSKEY \
-e TLSCERT=$TLSCERT \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e SESSIONKEY=$SESSIONKEY \
-e MESSAGE_ADDR=$MESSAGE_ADDR \
-e SUMMARY_ADDR=$SUMMARY_ADDR \
-e RABBIT=$RABBIT \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--name gateway \
nehay100/gateway

