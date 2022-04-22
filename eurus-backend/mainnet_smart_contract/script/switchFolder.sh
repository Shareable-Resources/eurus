#!/bin/sh
SCRIPT_PATH=$(dirname $0)
if [ $# -lt 1 ]; then
    echo "Usage: $(basename $0) [with_proxy | no_proxy]"
    exit 0
fi

cd $SCRIPT_PATH/../

rm  migrations
if [ "$1" = "with_proxy" ];then
    ln -s migrations_mainnet migrations
elif [ "$1" = "no_proxy" ]; then
    ln -s migrations_mainnet_noProxy migrations
else
    ln -s migrations_eurus_sidechain migrations
fi
