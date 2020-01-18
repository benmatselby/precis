# Precis

[![Build Status](https://travis-ci.org/benmatselby/precis.png?branch=master)](https://travis-ci.org/benmatselby/precis)

CLI application for getting information out of various systems in a neat dashboard format. This is inspired from [@jessfraz/tdash](https://github.com/jessfraz/tdash)'s dashboard

It integrates with:

- GitHub
- Jenkins

## Usage

```shell
.______   .______       _______   ______  __       _______.
|   _  \  |   _  \     |   ____| /      ||  |     /       |
|  |_)  | |  |_)  |    |  |__   |  ,----'|  |    |   (----
|   ___/  |      /     |   __|  |  |     |  |     \   \
|  |      |  |\  \----.|  |____ |   ----.|  | .----)   |
| _|      | _|  ._____||_______| \______||__| |_______/

A terminal dashboard which gives an overview of useful things

Build:

  -display-build
        Do you want to show build information from Jenkins? (default true)
  -display-github
    	Do you want to show GitHub information? (default true)
  -github-owner string
    	The GitHub CI owner (or define env var GITHUB_OWNER)
  -github-token string
    	The GitHub CI authentication token (or define env var GITHUB_TOKEN)
  -interval string
    	The refresh rate for the dashboard (default "60s")
  -jenkins-password string
    	The Jenkins password to authenticate with (or define env var JENKINS_PASSWORD)
  -jenkins-url string
    	The Jenkins URL (or define env var JENKINS_URL)
  -jenkins-username string
    	The Jenkins username to authenticate with (or define env var JENKINS_USERNAME)
  -jenkins-view string
    	The Jenkins view you want render, otherwise it is all (or define env var JENKINS_VIEW)
```

## Requirements

If you are wanting to build and develop this, you will need the following items installed.

- Go version 1.11+

## Configuration

You will need the following environment variables defining, depending on which systems you are running in the dashboard:

```shell
export GITHUB_TOKEN=""
export GITHUB_OWNER=""
export JENKINS_URL=""
export JENKINS_USERNAME=""
export JENKINS_PASSWORD=""
export JENKINS_VIEW=""
```

You can also define `~/.benmatselby/precis.yml` which has various settings.

### Limiting the repos to show Pull Requests for

```shell
github:
  pull_request_repos:
  - my-org/my-repo
  - benmatselby/*
```

## Installation via Git

```shell
git clone git@github.com:benmatselby/precis.git
cd precis
make all
./precis
```

You can also install into your `$GOPATH/bin` by `go install`
