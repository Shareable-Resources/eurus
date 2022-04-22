#!/bin/sh

if [ $# -eq 0 ] || [ "$#" = "" ]; then
  echo 'Please input sh run_server.sh { all | merchantDemo }'
  exit 1
fi

RUN_FOLDER=/home/ubuntu/merchant_demo

runMerchantDemo() {
  cd $RUN_FOLDER
    tmux new-session -d -s merchantDemoServer
    tmux new-window -t merchantDemoServer -n merchantDemoServer_Win "bash"
    tmux send-keys -t merchantDemoServer "node dist/merchant_admin_service/server/merchant_admin_service/index.js" ENTER
}

if [ "$1" = "merchantDemo" ]; then
    runMerchantDemo 
fi
