FROM debian:latest

ARG KERNEL_GIT=git://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git
ARG CONTAINERD_VERSION=1.6.9
ARG CONTAINERD_PKG=containerd-${CONTAINERD_VERSION}-linux-amd64.tar.gz
ARG CONTAINERD_URL=https://github.com/containerd/containerd/releases/download/v${CONTAINERD_VERSION}/${CONTAINERD_PKG}

ARG controller=ebpf-attachment-controller
ARG workdir=/opt/controller/
ARG controller_config=controller-config.json

ENV CONTROLLER_CONFIG ${workdir}/config/${controller_config}
ENV CONTROLLER_TOOLS ${workdir}/src/tools/

# Install dependencies
RUN apt-get update -y
RUN apt-get install -y make gcc libssl-dev bc libelf-dev libcap-dev \
	clang gcc-multilib llvm libncurses5-dev git pkg-config libmnl-dev \
	bison flex libbpf-dev iproute2 jq wget apt binutils-dev

# Install bpftool
RUN git clone --depth 1 -b master ${KERNEL_GIT} /tmp/linux && \
cd /tmp/linux/tools/bpf/bpftool/ &&\
sed -i '/CFLAGS += -O2/a CFLAGS += -static' Makefile && \
sed -i 's/LIBS = -lelf $(LIBBPF)/LIBS = -lelf -lz $(LIBBPF)/g' Makefile && \
printf 'feature-libbfd=0\nfeature-libelf=1\nfeature-bpf=1\nfeature-libelf-mmap=1' >> FEATURES_DUMP.bpftool && \
FEATURES_DUMP=`pwd`/FEATURES_DUMP.bpftool make -j `getconf _NPROCESSORS_ONLN` && \
strip bpftool && \
ldd bpftool 2>&1 | grep -q -e "Not a valid dynamic program" \
	-e "not a dynamic executable" || \
	( echo "Error: bpftool is not statically linked"; false ) && \
mv bpftool /usr/bin && rm -rf /tmp/linux

#Install containerd
RUN wget ${CONTAINERD_URL} && \
mv ${CONTAINERD_PKG} /tmp/ && \
tar -xvzf /tmp/${CONTAINERD_PKG} -C /tmp/ && \
mv /tmp/bin/* /usr/local/bin/ && \
rm -rf /tmp/${CONTAINERD_PKG}

COPY ./bin/${controller} ${workdir}/bin/${controller}

COPY ./bpf-attachments ${workdir}/bpf-attachments

COPY ./config/${controller_config} ${workdir}/config/${controller_config}

COPY ./src/tools ${workdir}/src/tools

RUN chmod a+x -R ${workdir}/src/tools

WORKDIR ${workdir}

ENTRYPOINT ["${workdir}/bin/${controller}"]