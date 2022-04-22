#!/bin/sh
SCRIPT_PATH=$(dirname $0)
if [ $# -lt 1 ]; then
    echo "Usage: $(basename $0) [dev/local/rinkeby/mainnet/dev_noproxy]"
    exit 0
fi

cd $SCRIPT_PATH/../

rm  migrations
if [ "$1" = "dev" ];then 
    ln -s migrations_eurus_sidechain migrations
elif [ "$1" = "dev_noproxy" ];then
     ln -s migrations_eurus_sidechain_withoutProxy migrations
elif [ "$1" = "rinkeby" ];then
    ln -s migrations_mainnet migrations
elif [ "$1" = "mainnet" ]; then
    ln -s migrations_mainnet migrations
else
    ln -s migrations_eurus_sidechain migrations
fi
