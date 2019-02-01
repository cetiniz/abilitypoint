#!/bin/bash

git clone https://github.com/Kitware/CMake/releases/download/v3.13.1/cmake-3.13.1-Linux-x86_64.sh /cmake-3.13.1-Linux-x86_64.sh
mkdir /opt/cmake
sh /cmake-3.13.1-Linux-x86_64.sh --prefix=/opt/cmake --skip-license
ln -s /opt/cmake/bin/cmake /usr/local/bin/cmake
cmake --version

git clone https://github.com/neo4j-drivers/seabolt.git ~/seabolt
~/seabolt/make_release.sh

OPENSSL_ROOT_DIR=/usr/local/opt/openssl
PKG_CONFIG_PATH=/root/seabolt/build/dist/share/pkgconfig
LD_LIBRARY_PATH=/root/seabolt/build/dist/lib
GOPATH=/go
