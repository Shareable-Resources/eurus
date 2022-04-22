#!/bin/sh

if [ "$1" != "" ]; then
    GOPATH=$1
fi

PKGPATH=$GOPATH/pkg/mod/github.com/ethereum/go-ethereum@v1.10.6
cd "$PKGPATH"
sudo chmod 777 $PKGPATH/build
env GO111MODULE=on go run build/ci.go install

if [ ! -f "$PKGPATH/build/bin/abigen" ]; then
    echo "Compile go ethereum failed"
else
    echo "Compile go ethereum success"
fi