#!/bin/sh

set -e

ret=0
tmux has-session -t=dapp-payment-api 2>/dev/null || ret=$?

if [ $ret -eq "0" ];
then
	tmux kill-window -t dapp-payment-api
fi

tmux new-session -d -s dapp-payment-api
tmux send-keys -t dapp-payment-api "./bin/debug/dapp-payment-api --config ./config.json" ENTER
