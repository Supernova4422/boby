#!/bin/bash

export tag=latest
export name=fld-bot-test
docker build --force-rm -t ${name}:${tag} -f dockerfile.test .
docker run --rm ${name}:${tag}
result=$?
docker rmi ${name}:${tag}
exit $result