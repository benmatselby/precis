# Precis

[![Build Status](https://travis-ci.org/benmatselby/precis.png?branch=master)](https://travis-ci.org/benmatselby/precis)

CLI application for getting information out of various systems in a neat dashboard format. This is inspired from [@jessfraz/tdash](https://github.com/jessfraz/tdash)'s dashboard

It integrates with:

- GitHub
- Azure DevOps
- TravisCI

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
  -display-azure-devops
        Do you want to show Azure DevOps information?
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
  -azure-devops-account string
        The Azure DevOps account (or define env var AZURE_DEVOPS_ACCOUNT)
  -azure-devops-build-branch string
        Comma separated list of branches to display (default "master")
  -azure-devops-build-count int
        How many builds should we display (default 10)
  -azure-devops-project string
        The Azure DevOps project (or define env var AZURE_DEVOPS_PROJECT)
  -azure-devops-team string
        The Azure DevOps team (or define env var AZURE_DEVOPS_TEAM)
  -azure-devops-token string
        The Azure DevOps auth token (or define env var AZURE_DEVOPS_TOKEN)
```

## Configuration

You will need the following environment variables defining, depending on which systems you are running in the dashboard:

```
$ export AZURE_DEVOPS_ACCOUNT=""
$ export AZURE_DEVOPS_PROJECT=""
$ export AZURE_DEVOPS_TEAM=""
$ export AZURE_DEVOPS_TOKEN=""
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
    -eAZURE_DEVOPS_ACCOUNT \
    -eAZURE_DEVOPS_PROJECT \
    -eAZURE_DEVOPS_TEAM \
    -eAZURE_DEVOPS_TOKEN \
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
