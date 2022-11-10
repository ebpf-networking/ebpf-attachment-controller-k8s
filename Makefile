BINDIR=bin
EXEC=ebpf-attachment-controller

BUILD_DIR=build
SRCDIR=src

PWD=$(shell pwd)

REGISTRY=docker.io
IMAGE=dushyantbehl/ebpf-attachment-controller
TAG=latest

.PHONY: all
.default: ${EXEC}

${EXEC}: ${BINDIR}
	cd ${SRCDIR} && go build -o ${PWD}/${BINDIR}/${EXEC}

${BINDIR}:
	mkdir -p ${BINDIR}

docker-build: ${EXEC}
	@docker build -t ${IMAGE}:${TAG} -f ${BUILD_DIR}/Dockerfile .
	@docker tag ${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:${TAG}

docker-push: docker-build
	@docker push ${IMAGE}:${TAG}

install:
#	@deploy/install.sh

uninstall:
#	@deploy/uninstall.sh

clean:
	@rm -rf ${BINDIR}
