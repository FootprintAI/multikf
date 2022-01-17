#!/usr/bin/env bash

# make sure the import packages are downloaded
#go mod vendor
./gomodtidy.sh

docker build -t footprintai/multikind:v1 \
 --no-cache -f Dockerfile .
docker push footprintai/multikind:v1
