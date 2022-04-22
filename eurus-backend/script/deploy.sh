OUTPUT_FOLDER=bin/release
RUN_FOLDER=$HOME/eurus/bin/

# For local test
#RUN_FOLDER=$HOME/Documents/GitHub/src/eurus-backend/bin/
if [ $# -eq 0 ] || [ "$1" = "" ] ; then
    echo 'Please input sh deploy.sh { all | config | auth | user | approval | withdraw | deposit | blockchain | userObs | background | sign | admin | merchantAdmin | kyc | sweep | serverControl | password | publicData | blockCypher }'
    exit 1
fi

if [ $1 = "config" ]; then
	tmux kill-session -t configServer
	cp $OUTPUT_FOLDER/configServer "$RUN_FOLDER"
  echo "configServer deployed"

elif [ $1 = "background" ]; then
  cp $OUTPUT_FOLDER/background "$RUN_FOLDER"
  echo "background deployed"

elif [ $1 = "walletBackground" ]; then
  cp $OUTPUT_FOLDER/walletBackground "$RUN_FOLDER"
  echo "walletBackground deployed"

elif [ $1 = "auth" ]; then
  tmux kill-session -t authServer
  cp $OUTPUT_FOLDER/authServer "$RUN_FOLDER"
  echo "authServer deployed"

elif [ $1 = "user" ]; then
  tmux kill-session -t userServer
  cp $OUTPUT_FOLDER/userServer "$RUN_FOLDER"
  echo "userServer deployed"

elif [ $1 = "kyc" ]; then
  tmux kill-session -t kycServer
  cp $OUTPUT_FOLDER/kycServer "$RUN_FOLDER"
  echo "kycServer deployed"

elif [ $1 = "sign" ]; then
  tmux kill-session -t signServer
  cp $OUTPUT_FOLDER/signServer "$RUN_FOLDER"
  echo "signServer deployed"

elif [ $1 = "userObs" ]; then
  tmux kill-window -t userObserver
  cp $OUTPUT_FOLDER/userObserver "$RUN_FOLDER"
  echo "userObserver deployed"
	exit 1

elif [ $1 = "admin" ]; then
  tmux kill-session -t adminServer
  cp $OUTPUT_FOLDER/adminServer "$RUN_FOLDER"
  echo "adminServer deployed"

elif [ $1 = "merchantAdmin" ]; then
  tmux kill-session -t merchantAdminServer
  cp $OUTPUT_FOLDER/merchantAdminServer "$RUN_FOLDER"
  echo "merchantAdminServer deployed"

elif [ $1 = "sweep" ]; then
  tmux kill-session -t sweepServer
  cp $OUTPUT_FOLDER/sweepServer "$RUN_FOLDER"
  echo "sweepServer deployed"

elif [ $1 = "blockchain" ]; then
    for i in $(seq 1 2);
    do
       tmux kill-session -t blockChainIndexer_"$i"
    done
  cp $OUTPUT_FOLDER/blockChainIndexer "$RUN_FOLDER"
    echo "blockChainIndexer deployed"

elif [ $1 = "withdraw" ]; then
  for i in $(seq 1 7);
  do
   tmux kill-session -t withdrawObserver_"$i"
  done
  cp $OUTPUT_FOLDER/withdrawObserver "$RUN_FOLDER"
  echo "withdrawObserver deployed"

elif [ $1 = "deposit" ]; then
  for i in $(seq 1 7);
  do
   tmux kill-session -t depositObserver_"$i"
  done
  cp $OUTPUT_FOLDER/depositObserver "$RUN_FOLDER"
   echo "depositObserver deployed"

elif [ $1 = "approval" ]; then
    for i in $(seq 1 4);
    do
      tmux kill-session -t approvalObserver_"$i"
    done
   cp $OUTPUT_FOLDER/approvalObserver "$RUN_FOLDER"
   echo "approvalObserver deployed"

elif [ $1 = "serverControl" ]; then
  cp $OUTPUT_FOLDER/serverControlClient "$RUN_FOLDER"
  echo "Server Control Client deployed"

elif [ $1 = "password" ]; then
  tmux kill-session -t passwordServer
  cp $OUTPUT_FOLDER/passwordServer "$RUN_FOLDER"
  echo "passwordServer deployed"
elif [ $1 = "publicData" ]; then
  tmux kill-session -t publicData
  cp -r typescript/dist/tsc/publicDataServer "$RUN_FOLDER"
  echo "publicDataServer deployed"
elif [ $1 = "blockCypher" ]; then
  cp $OUTPUT_FOLDER/blockCypherBackground "$RUN_FOLDER"
  echo "blockCypher deployed"
elif [ $1 = "all" ]; then
  tmux kill-session -t configServer
	tmux kill-session -t authServer
	tmux kill-session -t userServer
  tmux kill-session -t adminServer
  tmux kill-session -t merchantAdminServer
  tmux kill-session -t sweepServer
  tmux kill-session -t kycServer
  tmux kill-session -t signServer
  tmux kill-window -t userObserver
  tmux kill-session -t passwordServer

  #kill-session
    for i in $(seq 1 7);
    do
      tmux kill-session -t withdrawObserver_"$i"
      tmux kill-session -t depositObserver_"$i"
    done

    for i in $(seq 1 4);
    do
      tmux kill-session -t approvalObserver_"$i"
    done

    for i in $(seq 1 2);
    do
       tmux kill-session -t blockChainIndexer_"$i"
    done
  tmux kill-session -t publicData
  #kill-session end
 	
	cp $OUTPUT_FOLDER/configServer "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/authServer "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/userServer "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/blockChainIndexer "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/approvalObserver "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/depositObserver "$RUN_FOLDER"
	cp $OUTPUT_FOLDER/withdrawObserver "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/background "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/walletBackground "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/adminServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/merchantAdminServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/sweepServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/kycServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/signServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/userObserver "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/serverControlClient "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/passwordServer "$RUN_FOLDER"
  cp $OUTPUT_FOLDER/blockCypherBackground "$RUN_FOLDER"
  cp -r typescript/dist/tsc/publicDataServer "$RUN_FOLDER"
else 
	echo "Invalid argument"
	exit 1
fi;

exit 0
