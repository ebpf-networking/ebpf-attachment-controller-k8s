FROM debian:latest

ARG workdir=/opt/workloads

RUN apt-get update && \
    apt-get install -y \
    zip \
    unzip \
    libpcap-dev \
    bridge-utils \
    net-tools \
    iptables \
    python3 \
    tcpdump \
    build-essential \
    python3-pip \ 
    python3-scapy \
    iproute2

RUN mkdir -p ${workdir}

RUN mkdir -p ${workdir}/bpf-filter/
RUN mkdir -p ${workdir}/ebpf-ratelimiter/

COPY ./bpf-attachments/bpf-filter/loadgen/ ${workdir}/bpf-filter/
COPY ./bpf-attachments/ebpf-ratelimiter/loadgen/ ${workdir}/ebpf-ratelimiter/

WORKDIR ${workdir}

# loop
CMD ["tail", "-f", "/dev/null"]