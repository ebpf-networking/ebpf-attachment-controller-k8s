BINDIR=bin
EXEC=ebpf-attachment-controller

BUILD_DIR=build
SRCDIR=src

PWD=$(shell pwd)

ATTACHMENT_FOLDER=bpf-attachments
ATTACHMENTS=bpf-filter ebpf-ratelimiter

REGISTRY=docker.io

IMAGE=dushyantbehl/ebpf-attachment-controller
TAG=latest

LOADGEN_IMAGE=dushyantbehl/scapy-loadgen
LOADGEN_FILE=${ATTACHMENT_FOLDER}/loadgen.Dockerfile

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

docker-build-loadgen:
	@docker build -t ${LOADGEN_IMAGE}:${TAG} -f ${LOADGEN_FILE} .
	@docker tag ${LOADGEN_IMAGE}:${TAG} ${REGISTRY}/${LOADGEN_IMAGE}:${TAG}

docker-push-loadgen:
	@docker push ${LOADGEN_IMAGE}:${TAG}

${ATTACHMENTS}:
	make -C ${ATTACHMENT_FOLDER}/$@

install:
#	@deploy/install.sh

uninstall:
#	@deploy/uninstall.sh

clean:
	@rm -rf ${BINDIR}
