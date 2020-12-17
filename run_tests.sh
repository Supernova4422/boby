#!/bin/sh
tag=fld-bot-test
docker build --force-rm -t ${tag} . -f dockerfile.test
docker run --rm ${tag} 
docker rmi ${tag}