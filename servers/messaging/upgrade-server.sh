export MYSQL_ROOT_PASSWORD="password"
# export MYSQL_CONTAINER_NAME="mysqldemo"
export MYSQL_DB="userDB"
export MYSQL_ADDR="mysqldemo"
# export DSN="root:$MYSQL_ROOT_PASSWORD@tcp($MYSQL_CONTAINER_NAME:3306)/$MYSQL_DB"
export PRIVNETNAME="privnet"
export RABBIT="rabbit"

docker rm -f node

# ADDING MAIN.GO CODE
docker pull nehay100/nodecontainer

docker run -d \
--network $PRIVNETNAME \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_ADDR=$MYSQL_ADDR \
-e MYSQL_DB=$MYSQL_DB \
-e ADDR="node:80" \
-e RABBIT=$RABBIT \
--name node \
nehay100/nodecontainer