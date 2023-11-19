ARG MYSQL_VERSION=8.0.28

FROM mysql:$MYSQL_VERSION-debian

RUN apt-get update \
    && apt-get install -y wget curl sudo openssh-server git vim telnet iputils-ping net-tools iproute2 libdbd-mysql-perl

COPY ./authorized_keys /root/.ssh/authorized_keys
COPY ./id_rsa /root/.ssh/id_rsa
COPY ./id_rsa.pub /root/.ssh/id_rsa.pub
COPY ./.mylogin.cnf /root/.mylogin.cnf

RUN mkdir /run/sshd \
    && echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config \
    && echo '    StrictHostKeyChecking no' >> /etc/ssh/ssh_config \
    && chmod 700 /root/.ssh \
    && chmod 600 /root/.ssh/authorized_keys \
    && chmod 600 /root/.ssh/id_rsa \
    && chmod 600 /root/.ssh/.mylogin.cnf \
    && echo 'root:123456' | chpasswd \
    && echo "alias ll='ls -l'" >> /root/.bashrc

ARG XTRABACKUP_VERSION=8.0.28-21
ARG GLIBC=2.17

RUN curl -O https://repo.percona.com/apt/percona-release_latest.generic_all.deb \
    && apt-get install gnupg2 lsb-release ./percona-release_latest.generic_all.deb \
    && percona-release enable tools \
    && apt-get update \
    && apt-get install qpress

RUN wget https://downloads.percona.com/downloads/Percona-XtraBackup-8.0/Percona-XtraBackup-$XTRABACKUP_VERSION/binary/tarball/percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal.tar.gz \
    && tar -xvf percona-xtrabackup-${XTRABACKUP_VERSION}-Linux-x86_64.glibc$GLIBC-minimal.tar.gz \
    && mv percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal /usr/local/xtrabackup \
    && rm -rf percona-xtrabackup-$XTRABACKUP_VERSION-Linux-x86_64.glibc$GLIBC-minimal.tar.gz

ENV PATH=$PATH:/usr/local/xtrabackup/bin

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ENV MYSQL_MAJOR 8.0

VOLUME /var/lib/mysql

ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 3306 33060 22
CMD ["mysqld"]
