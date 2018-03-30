#!/bin/bash
set -ex
source "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/docker-tag.sh"
docker build -t ewilde/$(basename $2):$(git describe) --build-arg REPONAME=$(basename $1) --build-arg APPNAME=$(basename $2) --build-arg TAGS={version:\"$(git describe)\"} -f build/package/docker/$(basename $2)/Dockerfile .
docker tag ewilde/$(basename $2):$(git describe) ewilde/$(basename $2):$TAG
