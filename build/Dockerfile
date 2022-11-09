FROM golang:1.19.2

RUN apt-get update -y
RUN apt-get install -y make gcc libssl-dev bc libelf-dev libcap-dev \
	clang gcc-multilib llvm libncurses5-dev git pkg-config libmnl-dev \
	bison flex libbpf-dev iproute2

ARG controller=ebpf-attachment-controller
ARG workdir=/opt/controller/
ARG controller_config=controller-config.json

COPY ./bin/${controller} ${workdir}/bin/${controller}

COPY ./bpf-attachments ${workdir}/bpf-attachments

COPY ./config/${controller_config} ${workdir}/config/${controller_config}

COPY ./src/tools ${workdir}/src/tools

RUN chmod a+x -R ${workdir}/src/tools

#RUN git clone https://git.kernel.org/pub/scm/network/iproute2/iproute2.git ;\
#cd iproute2 ;\
#./configure --prefix='/usr' ;\
#make install ;\
# install -m 0644 ./include/bpf_api.h /usr/include/iproute2

#RUN cd ${workdir}/bpf-attachments/bpf-filter; make

WORKDIR ${workdir}

ENTRYPOINT ["./bin/${controller}"]