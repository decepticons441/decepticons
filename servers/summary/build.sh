# !/bin/bash
GOOS=linux go build
docker build -t nehay100/summary .
go clean

