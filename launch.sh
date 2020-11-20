#!/usr/bin/env bash

git pull
cd src/main
go build -o fld-bot
pkill fld-bot
./fld-bot &
