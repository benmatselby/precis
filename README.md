# Precis

[![Build Status](https://travis-ci.org/benmatselby/precis.png?branch=master)](https://travis-ci.org/benmatselby/precis)

CLI application for getting information out of various systems in a neat dashboard format. This is inspired from [@jessfraz/tdash](https://github.com/jessfraz/tdash)'s dashboard

## Usage

```
.______   .______       _______   ______  __       _______.
|   _  \  |   _  \     |   ____| /      ||  |     /       |
|  |_)  | |  |_)  |    |  |__   |  ,----'|  |    |   (----
|   ___/  |      /     |   __|  |  |     |  |     \   \
|  |      |  |\  \----.|  |____ |   ----.|  | .----)   |
| _|      | _|  ._____||_______| \______||__| |_______/

A terminal dashboard which gives an overview of useful things

Build:

  -current-iteration string
    	What is the current iteration
  -display-github
    	Do you want to show GitHub information? (default true)
  -display-travis
    	Do you want to show Travis CI information? (default true)
  -display-vsts
    	Do you want to show Visual Studio Team Services information?
  -github-owner string
    	The GitHub CI owner (or define env var GITHUB_OWNER)
  -github-token string
    	The GitHub CI authentication token (or define env var GITHUB_TOKEN)
  -interval string
    	The refresh rate for the dashboard (default "60s")
  -travis-owner string
    	The Travis CI owner (or define env var TRAVIS_CI_OWNER)
  -travis-token string
    	The Travis CI authentication token (or define env var TRAVIS_CI_TOKEN)
  -vsts-account string
    	The Visual Studio Team Services account (or define env var VSTS_ACCOUNT)
  -vsts-build-branch string
    	Comma separated list of branches to display (default "master")
  -vsts-build-count int
    	How many builds should we display (default 10)
  -vsts-project string
    	The Visual Studio Team Services project (or define env var VSTS_PROJECT)
  -vsts-team string
    	The Visual Studio Team Services team (or define env var VSTS_TEAM)
  -vsts-token string
    	The Visual Studio Team Services auth token (or define env var VSTS_TOKEN)
```

## Configuration

You will need the following environment variables defining, depending on which systems you are running in the dashboard:

```
$ export VSTS_ACCOUNT=""
$ export VSTS_PROJECT=""
$ export VSTS_TEAM=""
$ export VSTS_TOKEN=""
$ export TRAVIS_CI_TOKEN=""
$ export TRAVIS_CI_OWNER=""
$ export GITHUB_TOKEN=""
$ export GITHUB_OWNER=""
```

## Installation via Docker

Other than requiring [docker](http://docker.com) to be installed, there are no other requirements to run the application this way. This is the preferred method of running the `precis`. The image is [here](https://hub.docker.com/r/benmatselby/precis/).

```
$ docker run \
    --rm \
    -t \
    -eVSTS_ACCOUNT \
    -eVSTS_PROJECT \
    -eVSTS_TEAM \
    -eVSTS_TOKEN \
    -eTRAVIS_CI_TOKEN \
    -eTRAVIS_CI_OWNER \
    -eGITHUB_TOKEN \
    -eGITHUB_OWNER \
    benmatselby/precis
```

## Installation via Git

```
$ git clone git@github.com:benmatselby/precis.git
$ cd precis
$ make all
$ ./precis
```
