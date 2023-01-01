#!/bin/bash

# rm
sudo docker stop auto-deploy-build-env
sudo docker rm auto-deploy-build-env
sudo docker rmi org/auto-deploy-build-env:latest


# Debian apt源
# https://mirrors.tuna.tsinghua.edu.cn/help/debian/
# https
cat>sources_https.list<<EOF
# 选择你的 Debian 版本: buster
# 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释
deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free
deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free

deb https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free

deb https://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free
# deb-src https://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free
EOF
# http
cat>sources_http.list<<EOF
# 选择你的 Debian 版本: buster
# 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free
# deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free
deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free
# deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free

deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free
# deb-src http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free

deb http://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free
# deb-src http://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free
EOF


# Dockerfile
cat>Dockerfile<<EOF
# https://hub.docker.com/
# https://hub.docker.com/_/debian
FROM debian:buster

# 以root执行
USER root

# 设置镜像的时间格式与时区
#RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 支持 ll 命令
#RUN echo "alias ll='ls -l'" >> ~/.bashrc && source ~/.bashrc
RUN echo "alias ll='ls -l'" >> ~/.bashrc

# Debian apt源
COPY sources_http.list /etc/apt/sources.list
EOF


# docker build --tag, -t
# 镜像的名字及标签，通常 name:tag 或者 name 格式。
sudo docker build -f Dockerfile -t org/auto-deploy-build-env:latest .


# docker run
# https://docs.docker.com/engine/reference/commandline/run/
# -i: 以交互模式运行容器，通常与 -t 同时使用；
# -d: 后台运行容器，并返回容器ID；
# --name [name] 为容器指定一个名称，后续可以通过名字进行容器管理
# -p [host port]:[container port]，指定端口映射，格式为: 主机(宿主)端口:容器端口
sudo docker run \
-id \
--name auto-deploy-build-env \
-p 18022:22 \
-t org/auto-deploy-build-env:latest


# exec.sh
cat>exec.sh<<EOF
#!/bin/bash
# \$ docker inspect [CONTAINER ID]
# \$ cd [\${MergedDir}]

# 以root用户进入docker容器
# \$ docker exec -it -u root [CONTAINER ID | CONTAINER NAME] /bin/bash
sudo docker exec -it -u root auto-deploy-build-env /bin/bash
EOF
chmod +x exec.sh


# debian
:<<!
$ apt update
#$ apt upgrade

# vim
$ apt install vim

# ping
$ apt install -y iputils-ping

# telnet
$ apt install telnet

# ifconfig
$ apt install net-tools

// ssh server
$ apt install openssh-server
$ vim /etc/ssh/sshd_config
#PermitRootLogin prohibit-password
PermitRootLogin yes
#UsePAM yes
UsePAM no
$ service ssh restart
$ whoami
root
$ id
uid=0(root) gid=0(root) groups=0(root)
$ passwd

!