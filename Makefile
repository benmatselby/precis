NAME := precis
DOCKER_PREFIX = benmatselby

.DEFAULT_GOAL := explain
.PHONY: explain
explain: ## Explain what can be done
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
	### Targets
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

GITCOMMIT := $(shell git rev-parse --short HEAD)

.PHONY: clean
clean: ## Clean the local dependencies
	rm -fr vendor

.PHONY: install
install: ## Install the local dependencies
	go get ./...

.PHONY: vet
vet: ## Vet the code
	go vet ./...

.PHONY: lint
lint: ## Lint the code
	golint -set_exit_status $(shell go list ./...)

.PHONY: build
build: ## Build the application
	go build .

.PHONY: static
static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/$(NAME)/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=profile.out
	# go tool cover -func=coverage.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=profile.out

.PHONY: all
all: clean install lint vet build test ## Run everything

.PHONY: static-all
static-all: clean install vet static test ## Run everything
