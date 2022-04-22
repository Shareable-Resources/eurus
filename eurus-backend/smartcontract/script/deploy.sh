 #!/bin/sh
SCRIPT_PATH=$(dirname $0)

if [ $# -lt 1 ]; then
    echo "Usage: $(basename $0) [dev/local/rinkeby/mainnet]"
    exit 0
fi

sh switchFolder.sh $1

cd $SCRIPT_PATH/../

if [ "$1" == "dev" ];then 
    truffle migrate --network besu_dev
elif [ "$1" == "rinkeby" ];then
    truffle migrate --network rinkeby --f 8
elif [ "$1" == "mainnet" ]; then
    truffle migrate --network mainnet
else
    truffle migrate
fi
