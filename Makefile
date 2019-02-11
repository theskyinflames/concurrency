GOPATH ?= $(HOME)/go
GOBIN ?= $(HOME)/bin

VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

DOCKER_IMAGE_CSV = csvreader
DOCKER_IMAGE_CRM = crmintegrator

all: help

help:
	@echo
	@echo "VERSION: $(VERSION)"
	@echo "BRANCH: $(BRANCH)"
	@echo
	@echo "usage: make <command>"
	@echo
	@echo "commands:"
	@echo "    mod       		- populate vendor/ without updating it first"
	@echo "    build  	     	- build apps and installs them in $(GOBIN)"
	@echo "    test      		- run unit tests"
	@echo "    coverage  		- run unit tests and show coaverage on browser"
	@echo "    clean     		- remove generated files and directories"
	@echo "    run-crm     		- start the crm integrator locally, NOT AS A DOCKER CONAINER"
	@echo "    run-csv     		- start the csv reader locally, NOT AS A DOCKER CONAINER"
	@echo "    docker-build-crm	- build the crm integrator docker image"
	@echo "    docker-build-csv	- build the csv integrator docker image"	
	@echo "    docker-run-crm 	- use docker-compose to run the crm integrator docker image"
	@echo "    docker-run-csv 	- use docker-compose to run the csv reader docker image"
	@echo
	@echo "GOPATH: $(GOPATH)"
	@echo "GOBIN: $(GOBIN)"
	@echo

mod:
	@echo ">>> Populating vendor folder..."
	@go mod vendor

build:
	@echo ">>> Building app..."
	go install -v ./...
	@echo

test:
	@echo ">>> Running tests..."
	go test -count=1 -v ./...
	@echo

coverage:
	go test ./... -v -coverprofile=coverage.out && go tool cover -html=coverage.out

clean:
	@echo ">>> Cleaning..."
	go clean -i -r -cache -testcache
	@echo

run-crm:
	@echo ">>> Running ..."
	go run ./crmintegrator/main.go
	@echo

run-csv:
	@echo ">>> Running ..."
	go run ./csvreader/main.go
	@echo

docker-db:
	@echo ">>> Starting PostgreSQL container ..."
	docker-compose up db
	@echo

docker-build-csv:
	@echo ">>> Docker image building ..."
	docker build -f Dockerfile-csv -t $(DOCKER_IMAGE_CSV) .
	@echo

docker-build-crm:
	@echo ">>> Docker image building ..."
	docker build -f Dockerfile-crm -t $(DOCKER_IMAGE_CRM) .
	@echo

docker-run-csv:
	@echo ">>> Running Docker image ..."
	docker-compose up $(DOCKER_IMAGE_CSV)
	@echo

docker-run-crm:
	@echo ">>> Running Docker image ..."
	docker-compose up $(DOCKER_IMAGE_CRM)
	@echo
