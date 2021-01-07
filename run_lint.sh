#!/bin/bash

export tag=latest
export name=fld-bot-lint
docker build --force-rm -t ${name}:${tag} -f dockerfile.lint .
docker run --rm ${name}:${tag}
result=$?
docker rmi ${name}:${tag}
exit $result