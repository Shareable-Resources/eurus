#!/bin/sh

DIR=$(pwd)/$(dirname "$0")

echo "go.mod content"
cat $(pwd)/go.mod

echo "Downloading dependency packages"
go get -d ...
go get -d ./...

sh ${DIR}/setup_truffle.sh

TRUFFLE_PATH=${DIR}/../smartcontract/node_modules/.bin/

export PATH=$PATH:$TRUFFLE_PATH
cd ${DIR}


echo "Building tools"
#go build  -o ${DIR}/../bin/release/extractAbi ../tool/extract_abi
go build  -o ${DIR}/../bin/release/extractBin ../tool/extract_bin


echo "Building go ethereum"

sh build_go_ethereum.sh $GOPATH

cd ${DIR}/../typescript
npm init -y

cd ${DIR}/../

make compile_sc
make genabi

make compile_mainnet_sc
make gen_mainnet_abi