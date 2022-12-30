#!/bin/bash

CUR_DIR=$(cd $(dirname $0); pwd)
echo CUR_DIR: ${CUR_DIR}

OUTPUT_DIR=${CUR_DIR}/build
echo OUTPUT_DIR: ${OUTPUT_DIR}

# build
OUTPUT_PKG=${OUTPUT_DIR}/o
echo OUTPUT_PKG: ${OUTPUT_PKG}
cd ./src && go build -ldflags="-s -w" -o "${OUTPUT_PKG}"