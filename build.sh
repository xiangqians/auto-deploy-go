#!/bin/bash

curDir=$(cd $(dirname $0); pwd)
echo curDir: ${curDir}

outputDir=${curDir}/build
echo outputDir: ${outputDir}

# 删除 build 目录
rm -rf "${outputDir}"

# 创建 build 目录
mkdir -p "${outputDir}"

# cp
cp -r i18n "${outputDir}/"
cp -r static "${outputDir}/"
cp -r templates "${outputDir}/"
cp -r data "${outputDir}/"
cp -r script "${outputDir}/"

# go
os=`go env GOOS`
arch=`go env GOARCH`
pkgName=o_${os}_${arch}
outputPkg=${outputDir}/${pkgName}
echo outputPkg: ${outputPkg}
cd ./src && go build -ldflags="-s -w" -o "${outputPkg}"

# startup.sh
cat>${outputDir}/startup.sh<<EOF
# startup.sh
# \$ chmod +x ${pkgName} startup.sh
./${pkgName}
EOF
chmod +x "${outputDir}/startup.sh"