BINDIR=bin
EXEC=ebpf-attachment-controller

BUILD_DIR=build
SRCDIR=src

ATTACHMENT_FOLDER=bpf-attachments
ATTACHMENTS=bpf-filter ebpf-ratelimiter

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

docker-build: ${EXEC} ${ATTACHMENTS}
	@docker build -t ${IMAGE}:${TAG} -f ${BUILD_DIR}/Dockerfile .
	@docker tag ${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:${TAG}

docker-push: docker-build
	@docker push ${IMAGE}:${TAG}

${ATTACHMENTS}:
	make -C ${ATTACHMENT_FOLDER}/$@

install:
#	@deploy/install.sh

uninstall:
#	@deploy/uninstall.sh

clean:
	@rm -rf ${BINDIR}
