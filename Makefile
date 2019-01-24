export GO111MODULE := on
DOCKER_ABSOLUTE_PATH:=/usr/local/bin
DOCKER_FILE:=./build/docker/Dockerfile
MODULE_NAME:=$(shell sh -c 'cat go.mod | grep module | sed -e "s/module //"')
PROJECT_NAME:=challenge
ASSETS_PATH_NAME:=assets
FILE_NAME:=data
ENV= \
-e ASSETS_PATH_NAME=${DOCKER_ABSOLUTE_PATH}/${ASSETS_PATH_NAME} \
-e FILE_NAME=${FILE_NAME}
VOLUMES= \
-v ${PWD}/${ASSETS_PATH_NAME}:${DOCKER_ABSOLUTE_PATH}/${ASSETS_PATH_NAME}

all: prepare format
.PHONY: build
format:
	go fmt `go list ./... | grep -v /vendor/`
	goimports -w -local ${MODULE_NAME} `go list -f {{.Dir}} ./...`
prepare:
	go mod download
run: rm
	mkdir -p ${ASSETS_PATH_NAME}
	docker run --name ${PROJECT_NAME} --rm ${ENV} -p 50000:50000 -d ${PROJECT_NAME}
build: rm
	docker build --no-cache -f ${DOCKER_FILE} -t ${PROJECT_NAME} .
rm:
	-docker rm --force ${PROJECT_NAME}
