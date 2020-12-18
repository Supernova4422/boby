#!/bin/bash

# By appending "|| true" execution is allowed to continue even when 0 isn't returned.
export name=fld-bot
docker stop ${name} || true
docker rm ${name} || true

export tag=fld-bot
docker rmi ${tag} || true
docker build -t ${tag} -f dockerfile .

docker run --name ${name} -d ${tag}
