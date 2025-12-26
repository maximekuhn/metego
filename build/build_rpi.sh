#!/bin/bash

set -e

# Building go-sqlite3 for ARM requires additional compilation flags
# The easiest way to cross-compile it is to use Docker
docker build -f ./build/Dockerfile.build -t metego .
mkdir -p ./bin/rpi
id=$(docker run -d metego)
docker cp "$id":/metego/bin/web ./bin/rpi/web
docker stop "$id"
docker rm "$id"
