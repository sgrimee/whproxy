NS ?= quay.io
REPO = whproxy

NAME = $(REPO)

DOCKERCMD ?= docker

VERSION ?= $(BRANCH)

current_dir = $(shell pwd)

NO_CACHE ?= false

.PHONY: build test image push start stop rm release

default: build

test: 
	go get -t
	go test

build:
	go get
	go build -ldflags "-X main.version=`git describe --tags`"

build_linux:
	TODO

image: test build_linux
	$(DOCKERCMD) build --no-cache=$(NO_CACHE) --rm -t $(NS)/$(REPO):$(VERSION) .

push:
	$(DOCKERCMD) push $(NS)/$(REPO):$(VERSION)

release: build image push

start:
	$(DOCKERCMD) run -d --name $(NAME) $(PORTS) $(VOLUMES) $(ENV) $(NS)/$(REPO):$(VERSION)

stop:
	$(DOCKERCMD) stop $(NAME)

rm:
	$(DOCKERCMD) rm $(NAME)

