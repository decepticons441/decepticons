# !/bin/bash
# export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
# export MYSQL_CONTAINER_NAME="mysqldemo"
# export MYSQL_DB="userDB"
# export MYSQL_ADDR="$MYSQL_CONTAINER_NAME:3306"
# export DSN="root:$MYSQL_ROOT_PASSWORD@tcp($MYSQL_CONTAINER_NAME:3306)/$MYSQL_DB"
# export PRIVNETNAME="privnet"

GOOS=linux go build

# build the container image
docker build -t nehay100/nodecontainer .

go clean