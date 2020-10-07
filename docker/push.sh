#!/bin/bash
set -e
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

cd "$( dirname "${BASH_SOURCE[0]}" )"

for dir in */ ; do
    docker push m0rf30/pacur-${dir::-1}:$TRAVIS_TAG
done
