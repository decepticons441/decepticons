# !/bin/bash

docker build -t nehay100/mysqldemo .
go clean
docker push nehay100/mysqldemo
# docker rm -f mysqldemo

# # export MYSQL_ROOT_PASSWORD="some super-secret password"
# # generate 18 random bytes and base64 encode them
# export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

# docker run -d \
# --network host \
# # -p 3306:3306 \
# --name mysqldemo \
# -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
# -e MYSQL_DATABASE=userDB \
# nehay100/mysqldemo

# sleep 25

# docker run -it \
# --rm \
# --network host \
# nehay100/mysqldemo sh -c "mysql -h127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD"

# docker volume ls -qf dangling=true | xargs -r docker volume rm