[build]
mvn clean
mvn package

[target]
./target/test-1.0-SNAPSHOT.jar
#target/test-1.0-SNAPSHOT.jar.original
#./target/classes
#target/maven-archiver

[deploy]
#! /bin/bash

# Dockerfile
cat>Dockerfile<<EOF
# FROM命令定义构建镜像的基础镜像
FROM openjdk:13

# 以root执行
USER root

# 设置镜像的时间格式与时区
RUN cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 创建应用目录
RUN mkdir -p /opt/appl

# 拷贝
# 拷贝文件到容器
COPY ./test-1.0-SNAPSHOT.jar /opt/appl/
#COPY test-1.0-SNAPSHOT.jar.original /opt/appl/
# 拷贝目录到容器
#COPY ./classes /opt/appl/classes
#COPY maven-archiver /opt/appl/maven-archiver

# 暴露端口
#EXPOSE 8080

# WORKDIR：设置工作目录，即cd命令
WORKDIR /opt/appl

# 容器入口
# 设置JVM栈内存、最小堆内存、最大堆内存
ENTRYPOINT [ "java", "-Dfile.encoding=utf-8", "-Xss4096K", "-Xms256M", "-Xmx256M", "-jar", "/opt/appl/test-1.0-SNAPSHOT.jar" ]
EOF

# docker build
sudo docker build -f Dockerfile -t org/test:1.0 .

# docker run
sudo docker run \
-id \
--name test \
-p 8088:8080 \
-t org/test:1.0
