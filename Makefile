BINDIR=bin
EXEC=ebpf-attachment-controller

BUILDDIR=build
SRCDIR=src

PWD=$(shell pwd)

REGISTRY=
IMAGE=db/ebpf-attachment-controller
TAG=0.01

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
#	@docker push ${REGISTRY}/${IMAGE}:${TAG}

install:
#	@deploy/install.sh

uninstall:
#	@deploy/uninstall.sh

clean:
	@rm -rf ${BINDIR}
