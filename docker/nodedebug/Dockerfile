FROM ubuntu:focal
RUN DEBIAN_FRONTEND=noninteractive \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
            bcc \
            bpfcc-tools \
            curl \
            htop \
            iotop \
            iproute2 \
            jq \
            kmod \
            less \
            lm-sensors \
            nethogs \
            smartmontools \
            sysstat \
            tcpdump \
            tmux \
            tree \
            vim
CMD ["sleep", "50000"]
