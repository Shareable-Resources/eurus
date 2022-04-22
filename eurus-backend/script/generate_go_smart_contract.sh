#!/bin/sh


CURR_PATH=$(pwd)/$(dirname "$0")

if [ $# -lt 2 ]; then
    echo "Usage: `basename $0` <0 = sidechain/ 1 = mainnet> <input json path> <output file path>"
    exit 0
fi

if [ ! -f "$2" ]; then
    echo "Input file not found"
    exit 1
fi

INPUT_FILE_NAME=$(basename $2)
STRUCT_NAME=`echo "$INPUT_FILE_NAME" |  rev | cut -d. -f2 | rev`

SC_FOLDER_PREFIX=smartcontract
PKG_NAME=contract
if [ $1 -eq 1 ]; then
    SC_FOLDER_PREFIX=mainnet_smart_contract
    PKG_NAME=mainnet_contract
fi


ABI_FOLDER=$CURR_PATH/../$SC_FOLDER_PREFIX/build/abi/
BIN_FOLDER=$CURR_PATH/../$SC_FOLDER_PREFIX/build/bin/
JSON_FOLDER=$CURR_PATH/../$SC_FOLDER_PREFIX/build/contracts/

mkdir -p $ABI_FOLDER $BIN_FOLDER $3

$CURR_PATH/../bin/release/extractBin -f=$JSON_FOLDER -abi=$ABI_FOLDER -bin=$BIN_FOLDER

$GOPATH/pkg/mod/github.com/ethereum/go-ethereum@v1.10.6/build/bin/abigen  --abi=$ABI_FOLDER$STRUCT_NAME.abi --bin=$BIN_FOLDER$STRUCT_NAME.bin --pkg=$PKG_NAME --type $STRUCT_NAME > $3/$STRUCT_NAME.go
