FROM ubuntu:18.04
FROM golang:1.11 AS build

LABEL GOPATH=~/go
LABEL HOME=~

RUN apt-get update

RUN apt-get install -y openssl git build-essential libssl-dev

#install latest cmake
ADD https://github.com/Kitware/CMake/releases/download/v3.13.1/cmake-3.13.1-Linux-x86_64.sh /cmake-3.13.1-Linux-x86_64.sh
RUN mkdir /opt/cmake
RUN sh /cmake-3.13.1-Linux-x86_64.sh --prefix=/opt/cmake --skip-license
RUN ln -s /opt/cmake/bin/cmake /usr/local/bin/cmake
RUN cmake --version

ENV OPENSSL_ROOT_DIR=/usr/local/opt/openssl
ENV PKG_CONFIG_PATH=HOME/seabolt/build/dist/share/pkgconfig
ENV LD_LIBRARY_PATH=HOME/seabolt/build/dist/lib
ENV GOPATH=HOME/go

RUN git clone https://github.com/neo4j-drivers/seabolt.git ~/seabolt
RUN ~/seabolt/make_release.sh


RUN git clone https://github.com/cetiniz/abilitypoint.git HOME/go/src/

CMD ["go", "~/go/src/abilitypoint/server.go"]
