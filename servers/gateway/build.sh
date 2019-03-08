# !/bin/bash
GOOS=linux go build
docker build -t nehay100/gateway .
go clean

