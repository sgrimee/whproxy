NS ?= quay.io/sgrimee
REPO = whproxy

NAME = $(REPO)

DOCKERCMD ?= docker

VERSION ?= $(shell git describe --abbrev=0 HEAD --tags)

current_dir = $(shell pwd)

NO_CACHE ?= false

.PHONY: build test image push start stop rm release

default: build

test: 
	go test

build:
	go build -i -ldflags "-X main.version=$(VERSION)"

build_linux:
	env GOOS=linux GOARCH=amd64 go build -o $(NAME)_linux -i -ldflags "-X main.version=$(VERSION)"

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

