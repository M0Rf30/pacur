#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}" )"

for dir in */ ; do
    cd $dir
    sed -i -e "s|go get github.com/M0Rf30/pacur.*|go get github.com/M0Rf30/pacur # `date`|g" Dockerfile
    sudo docker build --rm -t pacur/${dir::-1} .
    sed -i -e "s|go get github.com/M0Rf30/pacur.*|go get github.com/M0Rf30/pacur|g" Dockerfile
    cd ..
done
