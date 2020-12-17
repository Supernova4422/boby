#!/bin/sh
name=fld-bot
docker stop ${name}
docker rm ${name}

tag=fld-bot
docker rmi ${tag}
docker build -t ${tag} . -f dockerfile

docker run --name ${name} -d ${tag}
