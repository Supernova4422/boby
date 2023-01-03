![Build](https://github.com/BKrajancic/boby/workflows/Build/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-81.6%25-brightgreen)
![Test](https://github.com/BKrajancic/boby/workflows/Test/badge.svg)
![Lint](https://github.com/BKrajancic/boby/workflows/golangci-lint/badge.svg)
![Vet](https://github.com/BKrajancic/boby/workflows/Vet/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/BKrajancic/boby)](https://goreportcard.com/report/github.com/BKrajancic/boby)

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-82%25-brightgreen.svg?longCache=true&style=flat)</a>

# boby
A configurable and flexible bot that can be used to make a unique bot with! What seperates two bots using this project is nothing more than some configuration files. 

# Servers this bot is used in
The primary implementation of this bot is known as FLD-Bot, which has been added to the following discord channels: 
1. Filipino learning and discussion
2. Tagalog.com

And more!

# Getting Started
To use to this project, the only required software is a working go environment. For installation instructions, see [this page.](https://golang.org/doc/install)

This project includes third party dependencies, so be sure to run `go get -d -v ./..` to install those dependencies.

This repository contains additional scripts and files that can be used to aid testing or running a bot. This includes a python3 script, and dockerfiles (which can be loaded by using the python3 scripts). These aren't essential, but can help with ensuring that your contributions work outside your own environment.

## Creating Configuration Files
To understand how to run the bot, first build it using `go build src/main`, then run the program to receive more instructions and examples on how to run the bot.

For properly understanding configuration files, make sure to view the files:

1. [goquery_scraper](https://github.com/BKrajancic/boby/blob/main/src/command/goquery_scraper.go)
2. [json_sender](https://github.com/BKrajancic/boby/blob/main/src/command/json_sender.go)
3. [regexp_scraper](https://github.com/BKrajancic/boby/blob/main/src/command/regexp_scraper.go)

Any of these files can be ignored by replacing its contents with `[]`.

Feel free to send a message if you are having issues running the bot. Unfortunately, this isn't an easy bot to configure.

##  Contributing
A pull request must have the following: 
1. New unit tests that are relavent to the commit
2. No failing unit tests
3. Documentation for everything with public scope

It would be desirable if your commit had the following:

1. golint returns no issues.
2. Tests coverage includes new and modified code. This repository is aiming for as high code coverage as possible, excluding the folders "service/discordservice" (because this code is  coupled to a 3rd party library, making testing difficult), "utils" and "main" (because they include side effects).  

## Adding bot to discord
To add your bot to a discord server with all the necesssary permissions, use the following
URL template:

> `https://discord.com/oauth2/authorize?client_id=<client_id>&permissions=0&scope=applications.commands%20bot`

## Logging TODOs and Issues
TODOs and issues are tracked using github's issue tracker.
