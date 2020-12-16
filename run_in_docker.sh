#!/bin/sh
tag=fld-bot
docker build -t ${tag} .
docker run -d ${tag}
