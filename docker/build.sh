#!/bin/bash
set -e
cd "$( dirname "${BASH_SOURCE[0]}" )"

if [ -z "${TRAVIS_TAG}"  ]; then
    TRAVIS_TAG="latest"
fi


for dir in */ ; do
    cd $dir
    sudo docker build --rm -t m0rf30/pacur-${dir::-1}:${TRAVIS_TAG} .
    cd ..
done
