#!/bin/sh
RUN_FOLDER=/home/ubuntu/eurus/bin/
# For local test
#RUN_FOLDER=$HOME/Documents/GitHub/src/eurus-backend/
PW=abcd1234
PW_SERVER_FLAG="--pwServer ${RUN_FOLDER}sock/pwServerSock"

adddate() {
  while IFS= read -r line; do
    printf '%s %s\n' "$(date)" "$line"
  done
}

if [ $# -eq 0 ] || [ "$#" = "" ]; then
  echo 'Please input sh run.sh { all | config | auth | user | approval | withdraw | deposit | blockchain | restart | sign | userObs | kyc | sweep | merchantAdmin | password | publicData | blockCypher }'
  exit 1
fi

runMerchantAdmin() {
  cd $RUN_FOLDER
  tmux new-session -d -s merchantAdminServer
  tmux new-window -t merchantAdminServer -n merchantAdminServer_Win "bash"
  tmux pipe-pane -t "merchantAdminServer:merchantAdminServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/MerchantAdminServer.log"
  tmux send-keys -t merchantAdminServer "./merchantAdminServer --config ../config/MerchantAdminServerConfig.json ${PW_SERVER_FLAG}" ENTER
}

runServerControlClient() {
  SOCK_PATH=$1
  if [ -f "../config/ServerControlClientConfig.json" ]; then
    ./serverControlClient --config ../config/ServerControlClientConfig.json ${PW_SERVER_FLAG}
  fi
}

runConfigServer() {
  tmux kill-window -t configServer
  cd $RUN_FOLDER
  tmux new-session -d -s configServer
  tmux new-window -t configServer -n configServer_Win "bash"
  tmux pipe-pane -t "configServer:configServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/ConfigServer.log"
  tmux send-keys -t configServer "./configServer --config ../config/ConfigServerConfig.json ${PW_SERVER_FLAG}" ENTER
  tmux send-keys -t configServer ${PW} ENTER


  sleep 5
}

runSignServer() {
  tmux kill-window -t signServer
  cd $RUN_FOLDER
  tmux new-session -d -s signServer
  tmux new-window -t signServer -n signServer_Win "bash"
  tmux pipe-pane -t "signServer:signServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/SignServer.log"
  tmux send-keys -t signServer "./signServer --config ../config/SignServerConfig.json ${PW_SERVER_FLAG}" ENTER
  tmux send-keys -t signServer ${PW} ENTER

}

runRestartService() {
  tmux kill-window -t restart
  tmux new-session -d -s restart
  tmux new-window -t restart -n restart_Win "bash"
  tmux pipe-pane -t "restart:restart_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/restart.log"
  tmux send-keys -t restart "sh $HOME/build/src/eurus-backend/script/restart.sh" ENTER
  sleep 5
}

runBackgroundService() {
  sh ./script/setup_crontab.sh
}

runWalletBackgroundService() {
  sh ./script/setup_crontab.sh
}

runBlockCypherBackgroundService() {
  sh ./script/setup_crontab.sh
}

runAuthServer() {
  tmux kill-window -t authServer
  cd $RUN_FOLDER
  tmux new-session -d -s authServer
  tmux new-window -t authServer -n authServer_Win "bash"
  tmux pipe-pane -t "authServer:authServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/AuthServer.log"
  tmux send-keys -t authServer "./authServer --config ../config/AuthServerConfig.json ${PW_SERVER_FLAG}" ENTER
  

  sleep 5
}

runUserServer() {
  tmux kill-window -t userServer
  cd $RUN_FOLDER
  tmux new-session -d -s userServer
  tmux new-window -t userServer -n userServer_Win "bash"
  tmux pipe-pane -t "userServer:userServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/UserServer.log"
  tmux send-keys -t userServer "./userServer --config ../config/UserServerConfig.json ${PW_SERVER_FLAG}" ENTER


  sleep 5
}

runKycServer() {
  tmux kill-window -t kycServer
  cd $RUN_FOLDER
  tmux new-session -d -s kycServer
  tmux new-window -t kycServer -n kycServer_Win "bash"
  tmux pipe-pane -t "kycServer:kycServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/KYCServer.log"
  tmux send-keys -t kycServer "./kycServer --config ../config/KYCServerConfig.json ${PW_SERVER_FLAG}" ENTER

  sleep 5
}

runSweepService() {
  tmux kill-window -t sweepServer
  cd $RUN_FOLDER
  tmux new-session -d -s sweepServer
  tmux new-window -t sweepServer -n sweepServer_Win "bash"
  tmux pipe-pane -t "sweepServer:sweepServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/SweepServer.log"
  tmux send-keys -t sweepServer "./sweepServer --config ../config/SweepServerConfig.json ${PW_SERVER_FLAG}" ENTER


  sleep 5
}

runUserObserver() {
  tmux kill-window -t userObserver
  cd $RUN_FOLDER
  tmux new-session -d -s userObserver
  tmux new-window -t userObserver -n userObserver_Win "bash"
  tmux pipe-pane -t "userObserver:userObserver_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/UserObserver.log"
  tmux send-keys -t userObserver "./userObserver --config ../config/UserObserverConfig.json ${PW_SERVER_FLAG}" ENTER
}

runPasswordServer() {
  tmux kill-window -t passwordServer
  cd $RUN_FOLDER
  tmux new-session -d -s passwordServer
  tmux new-window -t passwordServer -n passwordServer_Win "bash"
  tmux pipe-pane -t "passwordServer:passwordServer_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/PasswordServer.log"
  tmux send-keys -t passwordServer "./passwordServer --config ../config/PasswordServerConfig.json" ENTER
}

runBlockChainIndexer() {
  cd $RUN_FOLDER
  if [ -z "$2" ];then
    for i in $(seq 1 1); do
      tmux kill-session -t blockChainIndexer_"$i"
    done
  else
    tmux kill-session -t blockChainIndexer_"$2"
  fi
  if [ -z "$2" ];then
    
      tmux new-session -d -s blockChainIndexer_"$i"
      tmux new-window -t blockChainIndexer_"$i" -n blockChainIndexer_"$i"_Win "bash"
      tmux pipe-pane -t "blockChainIndexer_"$i":blockChainIndexer_"$i"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/BlockChainIndexer_"$i".log"
      tmux send-keys -t blockChainIndexer_"$i" "./blockChainIndexer --config ../config/BlockChainIndexerConfig_"$i".json  ${PW_SERVER_FLAG}" ENTER
    echo "blockChainIndexer started"
  else
    tmux new-session -d -s blockChainIndexer_"$2"
    tmux new-window -t blockChainIndexer_"$2" -n blockChainIndexer_"$2"_Win "bash"
    tmux pipe-pane -t "blockChainIndexer_"$2":blockChainIndexer_"$2"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/BlockChainIndexer_"$2".log"
    tmux send-keys -t blockChainIndexer_"$2" "$RUN_FOLDER./blockChainIndexer --config ../config/BlockChainIndexerConfig_"$2".json  ${PW_SERVER_FLAG}" ENTER
    echo "blockChainIndexer_${2} started"
  fi
}

runWithdrawObserver() {
  if [ -z "$2" ];then
    for i in $(seq 1 7); do
      tmux kill-window -t withdrawObserver_"$i"
    done
  else
    tmux kill-window -t withdrawObserver_"$2"
  fi
  cd $RUN_FOLDER
    if [ -z "$2" ];then
      for i in $(seq 1 7); do
        tmux new-session -d -s withdrawObserver_"$i"
        tmux new-window -t withdrawObserver_"$i" -n withdrawObserver_"$i"_Win "bash"
        tmux pipe-pane -t "withdrawObserver_"$i":withdrawObserver_"$i"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/WithdrawObserver_"$i".log"
        tmux send-keys -t withdrawObserver_"$i" "./withdrawObserver --config ../config/WithdrawObserverConfig_"$i".json ${PW_SERVER_FLAG}" ENTER
      done
       echo "7 withdrawObserver started"
    else
      
      tmux new-session -d -s withdrawObserver_"$2"
      tmux new-window -t withdrawObserver_"$2" -n withdrawObserver_"$2"_Win "bash"
      tmux pipe-pane -t "withdrawObserver_"$2":withdrawObserver_"$2"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/WithdrawObserver_"$2".log"
      tmux send-keys -t withdrawObserver_"$2" "$RUN_FOLDER./withdrawObserver --config ../config/WithdrawObserverConfig_"$2".json ${PW_SERVER_FLAG}" ENTER
      echo "withdrawObserver_$2 started"
    fi
}

runDepositObserver(){
   if [ -z "$2" ];then
  for i in $(seq 1 7); do
    tmux kill-window -t depositObserver_"$i"
    done
  else
     tmux kill-window -t depositObserver_"$2"
  fi
  
  cd $RUN_FOLDER
  if [ -z "$2" ];then
    for i in $(seq 1 7); do
      tmux new-session -d -s depositObserver_"$i"
      tmux new-window -t depositObserver_"$i" -n depositObserver_"$i"_Win "bash"
      tmux pipe-pane -t "depositObserver_"$i":depositObserver_"$i"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/DepositObserver_"$i".log"
      tmux send-keys -t depositObserver_"$i" "./depositObserver --config ../config/DepositObserverConfig_"$i".json ${PW_SERVER_FLAG}" ENTER
    done
      echo "7 depositObserver started"
  else
    
    tmux new-session -d -s depositObserver_"$2"
    tmux new-window -t depositObserver_"$2" -n depositObserver_"$2"_Win "bash"
    tmux pipe-pane -t "depositObserver_"$2":depositObserver_"$2"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/DepositObserver_"$2".log"
    tmux send-keys -t depositObserver_"$2" "$RUN_FOLDER./depositObserver --config ../config/DepositObserverConfig_"$2".json ${PW_SERVER_FLAG}" ENTER
    
    echo "depositObserver_$2 started"
  fi
}

runApprovalObserver() {
    if [ -z "$2" ];then
    for i in $(seq 1 2); do
      tmux kill-window -t approvalObserver_"$i"
    done
  else
    tmux kill-window -t approvalObserver_"$2"
  fi

  cd $RUN_FOLDER
  if [ -z "$2" ];then
    for i in $(seq 1 2); do
      tmux new-session -d -s approvalObserver_"$i"
      tmux new-window -t approvalObserver_"$i" -n approvalObserver_"$i"_Win "bash"
      tmux pipe-pane -t "approvalObserver_"$i":approvalObserver_"$i"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/ApprovalObserver_"$i".log"
      tmux send-keys -t approvalObserver_"$i" "./approvalObserver --config ../config/ApprovalObserver_"$i".json ${PW_SERVER_FLAG}" ENTER
    done
    echo "4 approvalObserver started"
  else
    
    tmux new-session -d -s approvalObserver_"$2"
    tmux new-window -t approvalObserver_"$2" -n approvalObserver_"$i"_Win "bash"
    tmux pipe-pane -t "approvalObserver_"$2":approvalObserver_"$2"_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/ApprovalObserver_"$2".log"
    tmux send-keys -t approvalObserver_"$2" "$RUN_FOLDER./approvalObserver --config ../config/ApprovalObserver_"$2".json ${PW_SERVER_FLAG}" ENTER
   echo "approvalObserver_$2 started"
  fi
}

runPublicData(){
  tmux kill-window -t publicData
  cd $RUN_FOLDER
  tmux new-session -d -s publicData
  tmux new-window -t publicData -n publicData_Win "bash"
  tmux pipe-pane -t "publicData:publicData_Win" "cat | ( while read line; do echo  $(date)  $"line"; done; ) >> ../errorLog/PublicDataServer.log"
  tmux send-keys -t publicData "node $RUN_FOLDER/publicDataServer/server/publicDataServer  $RUN_FOLDER/../config/PublicDataServerConfig.json" ENTER
}


if [ $1 = "config" ]; then
  runConfigServer
  exit 0

elif [ $1 = "sign" ]; then
  runSignServer
  exit 0

elif [ $1 = "restart" ]; then
  runRestartService
  exit 0

elif [ $1 = "background" ]; then
  runBackgroundService
  exit 0
elif [ $1 = "blockCypher" ]; then
  runBlockCypherBackgroundService
  exit 0
elif [ $1 = "walletBackground" ]; then
  runWalletBackgroundService
  exit 0
elif [ $1 = "auth" ]; then
  runAuthServer
  exit 0

elif [ $1 = "user" ]; then
  runUserServer
  exit 0

elif [ $1 = "kyc" ]; then
  runKycServer
  exit 0

elif [ $1 = "sweep" ]; then
  runSweepService
  exit 0

elif [ $1 = "userObs" ]; then
  runUserObserver
  exit 0

elif [ $1 = "blockchain" ]; then
  runBlockChainIndexer $1 $2
  exit 0

elif [ $1 = "withdraw" ]; then
  runWithdrawObserver $1 $2
  exit 0

elif [ $1 = "deposit" ]; then
  runDepositObserver $1 $2
  exit 0

elif [ $1 = "approval" ]; then
  runApprovalObserver $1 $2
  exit 0
elif [ $1 = "merchantAdmin" ]; then
  runMerchantAdmin
  exit 0
elif [ $1 = "password" ]; then
  runPasswordServer
  exit 0
 elif [ $1 = "publicData" ]; then
  runPublicData
  exit 0
elif [ $1 = "all" ]; then
  runPasswordServer
  runConfigServer
  runAuthServer
  runKycServer
  runUserObserver
  runBlockChainIndexer $1 $2
  runDepositObserver $1 $2
  runWithdrawObserver $1 $2
  runApprovalObserver $1 $2
  runSweepService
  runBackgroundService
  runWalletBackgroundService
  runBlockCypherBackgroundService
  runSignServer
  runUserObserver
  runMerchantAdmin
  runUserServer
  runRestartService
  runPublicData
else
  echo "Invalid argument"
  exit 1
fi
exit 0
