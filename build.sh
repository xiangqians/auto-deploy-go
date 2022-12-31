#!/bin/bash

curDir=$(cd $(dirname $0); pwd)
echo curDir: ${curDir}

outputDir=${curDir}/build
echo outputDir: ${outputDir}

# 删除 build 目录
rm -rf "${outputDir}"
echo rd: ${outputDir}

# 创建 build 目录
mkdir -p "${outputDir}"

# cp
cp -r i18n "${outputDir}/"
cp -r static "${outputDir}/"
cp -r templates "${outputDir}/"
cp -r data "${outputDir}/"
cp -r script "${outputDir}/"

# pkgName
os=`go env GOOS`
arch=`go env GOARCH`
pkgName=o_${os}_${arch}
echo pkgName: ${pkgName}

# go
pkgPath=${outputDir}/${pkgName}
#cd ./src && go build -ldflags="-s -w" -o "${pkgPath}"
# \$ sudo apt install upx
cd ./src && go build -ldflags="-s -w" -o "${pkgPath}" && upx -9 --brute "${pkgPath}"
echo pkgPath: ${pkgPath}

# startup.sh
startupPath=${outputDir}/startup.sh
cat>${startupPath}<<EOF
# startup.sh
# \$ chmod +x ${pkgName} startup.sh
./${pkgName}
EOF
echo startupPath: ${startupPath}
chmod +x "${startupPath}"