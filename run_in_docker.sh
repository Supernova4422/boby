#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "A single parameter must be passed, which becomes the container's name."
    exit 1
fi


# By appending "|| true" execution is allowed to continue even when 0 isn't returned.
docker stop ${1} || true
docker rm ${1} || true

export tag=latest
docker rmi ${1}:${tag} || true

docker build -t ${1}:${tag} -f dockerfile .
if [ $? != 0 ]
then 
    exit $? 
fi

docker run --name ${1} --restart always -d ${1}:${tag}
