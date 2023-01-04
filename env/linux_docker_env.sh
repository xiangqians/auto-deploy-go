#!/bin/bash

# rm
sudo docker stop auto-deploy-env
sudo docker rm auto-deploy-env
sudo docker rmi org/auto-deploy-env:latest


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
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 支持 ll 命令
#RUN echo "alias ll='ls -l'" >> ~/.bashrc && source ~/.bashrc
RUN echo "alias ll='ls -l'" >> ~/.bashrc

# Debian apt源
COPY sources_http.list /etc/apt/sources.list

# 创建目录
RUN mkdir -p /opt/auto-deploy

# 拷贝目录到容器
#COPY ./data /opt/auto-deploy/data
COPY ./i18n /opt/auto-deploy/i18n
COPY ./static /opt/auto-deploy/static
COPY ./templates /opt/auto-deploy/templates

# 拷贝文件到容器
COPY ./auto_deploy_linux_amd64 /opt/auto-deploy/
COPY ./startup.sh /opt/auto-deploy/

# 授予可执行权限
RUN chmod +x /opt/auto-deploy/auto_deploy_linux_amd64
RUN chmod +x /opt/auto-deploy/startup.sh


# 暴露端口
#EXPOSE 8080

# WORKDIR：设置工作目录，即cd命令
WORKDIR /opt/auto-deploy

# 容器入口
#ENTRYPOINT [ "./startup.sh" ]
EOF


# docker build --tag, -t
# 镜像的名字及标签，通常 name:tag 或者 name 格式。
sudo docker build -f Dockerfile -t org/auto-deploy-env:latest .


# docker run
# https://docs.docker.com/engine/reference/commandline/run/
# -i: 以交互模式运行容器，通常与 -t 同时使用；
# -d: 后台运行容器，并返回容器ID；
# --name [name] 为容器指定一个名称，后续可以通过名字进行容器管理
# -p [host port]:[container port]，指定端口映射，格式为: 主机(宿主)端口:容器端口
# -v [宿主机路径]:[容器路径]
sudo docker run \
-id \
--name auto-deploy-env \
-v /var/run/docker.sock:/var/run/docker.sock \
-v /usr/bin/docker:/usr/bin/docker \
-v /opt/auto-deploy/data:/opt/auto-deploy/data \
-p 18080:8080 \
-t org/auto-deploy-env:latest


# exec.sh
cat>exec.sh<<EOF
#!/bin/bash
# \$ docker inspect [CONTAINER ID]
# \$ cd [\${MergedDir}]

# 以root用户进入docker容器
# \$ docker exec -it -u root [CONTAINER ID | CONTAINER NAME] /bin/bash
sudo docker exec -it -u root auto-deploy-env /bin/bash
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

# ssh server
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

# ca
$ apt install --no-install-recommends ca-certificates curl

# 设置开机自启 auto-deploy 应用
$ cd /etc/init.d/
$ cat>autodeploy<<EOF
#!/bin/sh
### BEGIN INIT INFO
# Provides:          autodeploy
# Required-Start:    \$local_fs \$syslog \$network
# Required-Stop:     \$local_fs \$syslog \$network
# Default-Start:     1 2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Auto deploy.
# Description:       Auto deploy.
### END INIT INFO
cd /opt/auto-deploy && ./startup.sh
exit 0
EOF
# 授予脚本可执行权限
$ chmod +x autodeploy
# 将启动脚本加入开机启动项
$ update-rc.d autodeploy defaults
#$ runlevel
#$ ll /etc/rc1.d
#$ update-rc.d -f autodeploy remove


!