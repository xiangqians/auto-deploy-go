#!/bin/bash

# rm
sudo docker stop auto-deploy-build-env
sudo docker rm auto-deploy-build-env
sudo docker rmi org/auto-deploy-build-env:latest


# Debian apt源
cat>sources.list<<EOF
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
FROM debian:buster
USER root
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "alias ll='ls -l'" >> ~/.bashrc
COPY sources.list /etc/apt/sources.list
EOF


# docker build
sudo docker build -f Dockerfile -t org/auto-deploy-build-env:latest .


# docker run
sudo docker run \
-id \
--name auto-deploy-build-env \
-p 18022:22 \
-t org/auto-deploy-build-env:latest


# exec.sh
cat>exec.sh<<EOF
sudo docker exec -it -u root auto-deploy-build-env /bin/bash
EOF
chmod +x exec.sh