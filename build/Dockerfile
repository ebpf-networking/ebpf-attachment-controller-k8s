FROM golang:1.19.2

RUN apt-get update -y
RUN apt-get install -y make gcc libssl-dev bc libelf-dev libcap-dev \
	clang gcc-multilib llvm libncurses5-dev git pkg-config libmnl-dev \
	bison flex libbpf-dev iproute2 jq wget

RUN wget https://github.com/containerd/containerd/releases/download/v1.6.9/containerd-1.6.9-linux-amd64.tar.gz
RUN tar -xvzf containerd-1.6.9-linux-amd64.tar.gz
RUN mv bin/* /usr/bin/

ARG controller=ebpf-attachment-controller
ARG workdir=/opt/controller/
ARG controller_config=controller-config.json

COPY ./bin/${controller} ${workdir}/bin/${controller}

COPY ./bpf-attachments ${workdir}/bpf-attachments

COPY ./config/${controller_config} ${workdir}/config/${controller_config}

COPY ./src/tools ${workdir}/src/tools

RUN chmod a+x -R ${workdir}/src/tools

WORKDIR ${workdir}

ENTRYPOINT ["./bin/${controller}"]