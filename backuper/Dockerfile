ARG MYSQL_VERSION=8.0.28

FROM mysql:$MYSQL_VERSION-debian

ARG XTRABACKUP_VERSION=8.0.28-21
ARG GLIBC=2.17

RUN apt-get update \
    && apt-get install -y wget curl sudo openssh-server git vim telnet iputils-ping net-tools iproute2 libdbd-mysql-perl \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && curl -O https://repo.percona.com/apt/percona-release_latest.generic_all.deb \
    && apt-get install -y gnupg2 lsb-release ./percona-release_latest.generic_all.deb \
    && percona-release enable tools \
    && apt-get update \
    && apt-get install -y qpress lz4 zstd \
    && wget https://downloads.percona.com/downloads/Percona-XtraBackup-8.0/Percona-XtraBackup-$XTRABACKUP_VERSION/binary/tarball/percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal.tar.gz \
    && tar -xvf percona-xtrabackup-${XTRABACKUP_VERSION}-Linux-x86_64.glibc$GLIBC-minimal.tar.gz \
    && mv percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal /usr/local/xtrabackup \
    && rm -rf percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal.tar.gz

ENV PATH=$PATH:/usr/local/xtrabackup/bin

RUN apt-get install -y sshpass \
    && wget https://mirrors.aliyun.com/golang/go1.21.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

ARG PASSWORD='root:Y3+qoXfJj^'
COPY ./init.sh /usr/local/bin/init.sh

RUN mkdir /run/sshd \
    && echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config \
    && echo '    StrictHostKeyChecking no' >> /etc/ssh/ssh_config \
    && echo "alias ll='ls -l'" >> /root/.bashrc \
    && echo 'set mouse-=a' >> /usr/share/vim/vim81/defaults.vim \
    && ssh-keygen -t rsa -N "" -f ~/.ssh/id_rsa \
    && echo $PASSWORD | chpasswd \
    && chmod +x /usr/local/bin/init.sh

ENV GOROOT=/usr/local/go
ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /root

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
