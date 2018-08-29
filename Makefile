NAME := precis
DOCKER_PREFIX = benmatselby

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
	dep ensure

.PHONY: vet
vet: ## Vet the code
	go vet -v ./...

.PHONY: build
build: ## Build the application
	go build .

.PHONY: static
static: ## Build the application
	CGO_ENABLED=0 go build -ldflags "-extldflags -static -X github.com/benmatselby/$(NAME)/version.GITCOMMIT=$(GITCOMMIT)" -o $(NAME) .

.PHONY: test
test: ## Run the unit tests
	go test ./... -coverprofile=profile.out

.PHONY: test-cov
test-cov: test ## Run the unit tests with coverage
	go tool cover -html=profile.out

.PHONY: all
all: clean install vet build test ## Run everything

.PHONY: static-all
static-all: clean install vet static test ## Run everything
