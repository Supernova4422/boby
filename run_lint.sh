#!/bin/bash

export tag=latest
export name=boby-lint
docker build --force-rm -t ${name}:${tag} -f dockerfile.lint .
docker run --rm ${name}:${tag}
result=$?
docker rmi ${name}:${tag}
exit $result