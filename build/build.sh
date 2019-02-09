#!/bin/bash

CUR_DIR=$(pwd)

cd $GOPATH/src/github.com/cetiniz/abilitypoint/cmd/
export GOOS=linux
sudo go build -o abilitypoint-linux-amd64
export GOOS=darwin
cd $CUR_DIR
