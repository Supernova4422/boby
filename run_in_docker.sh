#!/bin/bash

# By appending "|| true" execution is allowed to continue even when 0 isn't returned.
export name=fld-bot
docker stop ${name} || true
docker rm ${name} || true

export tag=latest
docker rmi ${name}:${tag} || true

docker build -t ${name}:${tag} -f dockerfile .
if [ $? != 0 ]
then 
    exit $? 
fi

docker run --name ${name} -d ${name}:${tag}