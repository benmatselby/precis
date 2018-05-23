NAME := precis
DOCKER_PREFIX = benmatselby

.PHONY: explain
explain:
	### Welcome
	#
	# .______   .______       _______   ______  __       _______.
	# |   _  \  |   _  \     |   ____| /      ||  |     /       |
	# |  |_)  | |  |_)  |    |  |__   |  ,----'|  |    |   (----
	# |   ___/  |      /     |   __|  |  |     |  |     \   \
	# |  |      |  |\  \----.|  |____ |   ----.|  | .----)   |
	# | _|      | _|  ._____||_______| \______||__| |_______/
	#
	### Installation
	#
	# $$ make clean install
	#

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean:
	rm -fr vendor

.PHONY: install
install:
	dep ensure

.PHONY: vet
vet:
	go vet -v ./...

.PHONY: build
build:
	go build .

.PHONY: static
static:
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/$(NAME)/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test:
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test
	go tool cover -html=profile.out

.PHONY: all
all: clean install vet build test

.PHONY: static-all
static-all: clean install vet static test
