# rm
sudo docker stop auto-deploy-build-env
sudo docker rm auto-deploy-build-env
sudo docker rmi org/auto-deploy-build-env:latest

# Debian apt源
# https://mirrors.tuna.tsinghua.edu.cn/help/debian/
cat>sources.list<<EOF
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

# Dockerfile
cat>Dockerfile<<EOF
# https://hub.docker.com/
# https://hub.docker.com/_/debian
# https://hub.docker.com/_/debian/tags?page=1&name=stable
FROM debian:stable

# 以root执行
USER root

# 设置镜像的时间格式与时区
#RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 支持 ll 命令
#RUN echo "alias ll='ls -l'" >> ~/.bashrc && source ~/.bashrc
RUN echo "alias ll='ls -l'" >> ~/.bashrc

# Debian apt源
#COPY sources.list /etc/apt/sources.list
EOF

# docker build --tag, -t
# 镜像的名字及标签，通常 name:tag 或者 name 格式。
sudo docker build -f Dockerfile -t org/auto-deploy-build-env:latest .

# docker run
# https://docs.docker.com/engine/reference/commandline/run/
# -i: 以交互模式运行容器，通常与 -t 同时使用；
# -d: 后台运行容器，并返回容器ID；
# --name [name] 为容器指定一个名称，后续可以通过名字进行容器管理
sudo docker run \
-id \
--name auto-deploy-build-env \
-t org/auto-deploy-build-env:latest

# exec.sh
cat>exec.sh<<EOF
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

!