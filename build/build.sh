#!/bin/bash
export GOOS=linux
export GOOS=darwin

CUR_DIR=$(pwd)

cd $GOPATH/src/github.com/cetiniz/abilitypoint/cmd/
sudo go build -o abilitypoint-linux-amd64
cd $CUR_DIR
