#!/usr/bin/env bash

git pull
cd src/main
go build -o fld-bot
echo "Killing current processes"
pkill fld-bot
echo "Starting new process!"
./fld-bot &
