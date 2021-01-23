![Build](https://github.com/BKrajancic/FLD-Bot/workflows/Build/badge.svg)
![Test](https://github.com/BKrajancic/FLD-Bot/workflows/Test/badge.svg)
![Lint](https://github.com/BKrajancic/FLD-Bot/workflows/Lint/badge.svg)
![Vet](https://github.com/BKrajancic/FLD-Bot/workflows/Vet/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/BKrajancic/FLD-Bot/internal)](https://goreportcard.com/report/github.com/BKrajancic/FLD-Bot/internal)

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-96%25-brightgreen.svg?longCache=true&style=flat)</a>

# FLD-Bot
A configurable and flexible bot that can be used to make a unique bot with! What seperates two bots using this project is nothing more than some configuration files. 

# Servers this bot is used in
The primary implementation of this bot is known as FLD-Bot, which has been added to the following discord channels: 
1. Filipino learning and discussion
2. Tagalog.com

# Getting Started
To use to this project, the only required software is a working go environment. For installation instructions, see [this page.](https://golang.org/doc/install)

This project includes third party dependencies, so be sure to run `go get -d -v ./..` to install those dependencie.

This repository contains additional scripts and files that can be used to aid testing or running a bot. This includes a python3 script, and dockerfiles (which can be loaded by using the python3 scripts). These aren't essential, but can help with ensuring that your contributions work outside your own environment.

When running the bot, you will likely have issues that configuration files are missing. This is an issue, and a task is to have the bot guide a user to what files need to be created when running the bot.

##  Contributing
A pull request must have the following: 
1. New unit tests that are relavent to the commit
2. No failing unit tests
3. Documentation for everything with public scope

It would be desirable if your commit had the following:

1. golint returns no issues.
2. Tests coverage includes new and modified code. This repository is aiming for as high code coverage as possible, excluding the folders "service/discordservice" (because this code is  coupled to a 3rd party library, making testing difficult), "utils" and "main" (because they include side effects).  

## Logging TODOs and Issues
TODOs and issues are tracked using github's issue tracker.
