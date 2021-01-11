#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "A single parameter must be passed, which becomes the container's name."
    exit 1
fi

export tag=latest
export name=${1}
docker build --force-rm -t ${name}:${tag} -f dockerfile.test .
docker run --rm ${name}:${tag}
result=$?
docker rmi ${name}:${tag}
exit $result
