FROM debian:stable

# 以root执行
USER root

# 设置镜像的时间格式与时区
#RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
#RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 创建目录
RUN mkdir -p /opt/auto-deploy

# 将多个目录和文件 copy 到 /opt/auto-deploy 目录下
COPY i18n /opt/auto-deploy/i18n
COPY static /opt/auto-deploy/static
COPY templates /opt/auto-deploy/templates
COPY data /opt/auto-deploy/data
#COPY o_linux_amd64 /opt/auto-deploy/
COPY o_linux_amd64 /opt/auto-deploy/o

# 可执行文件
RUN chmod +x /opt/auto-deploy/o

# 暴露端口
#EXPOSE 8080

# WORKDIR：设置工作目录，即cd命令
WORKDIR /opt/auto-deploy

# 容器入口
ENTRYPOINT ["./o"]