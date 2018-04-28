Precis
======

[![Build Status](https://travis-ci.org/benmatselby/precis.png?branch=master)](https://travis-ci.org/benmatselby/precis)

CLI application for getting information out of various systems in a neat dashboard format. This is inspired from [@jessfraz/tdash](https://github.com/jessfraz/tdash)'s dashboard

# Usage

```
.______   .______       _______   ______  __       _______.
|   _  \  |   _  \     |   ____| /      ||  |     /       |
|  |_)  | |  |_)  |    |  |__   |  ,----'|  |    |   (----
|   ___/  |      /     |   __|  |  |     |  |     \   \
|  |      |  |\  \----.|  |____ |   ----.|  | .----)   |
| _|      | _|  ._____||_______| \______||__| |_______/

A terminal dashboard which gives an overview of useful things

Build:

  -d	Run in debug mode
  -interval string
    	The refresh rate for the dashboard (default "60s")
  -travis-owner string
    	The Travis CI owner (or define env var TRAVIS_CI_OWNER)
  -travis-token string
    	The Travis CI authentication token (or define env var TRAVIS_CI_TOKEN)
  -vsts-account string
    	The Visual Studio Team Services account (or define env var VSTS_ACCOUNT)
  -vsts-build-count int
    	How many builds should we display (default 10)
  -vsts-project string
    	The Visual Studio Team Services project (or define env var VSTS_PROJECT)
  -vsts-team string
    	The Visual Studio Team Services team (or define env var VSTS_TEAM)
  -vsts-token string
    	The Visual Studio Team Services auth token (or define env var VSTS_TOKEN)
```

# Installation via Docker

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
    benmatselby/precis
```


# Installation via Git

```
$ git clone git@github.com:benmatselby/precis.git
$ cd precis
$ make all
$ ./precis
```
