
if [ $# -eq 0 ] ; then
    echo 'Please input sh stop_observer.sh { all | config | auth | user | withdraw | deposit | approval | blockchain | restart | sweep }'
    exit 1
fi

if [ $1 = "config" ]; then
	tmux kill-session -t configServer
  echo "configServer stop"
	exit 1

elif [ $1 = "auth" ]; then
  tmux kill-session -t authServer
  echo "authServer stop"
	exit 1

elif [ $1 = "sign" ]; then
  tmux kill-session -t signServer
  echo "signServer stop"
	exit 1

elif [ $1 = "userObserver" ]; then
  tmux kill-window -t userObserver
  echo "userObserver stop"
	exit 1

elif [ $1 = "restart" ]; then
  tmux kill-session -t restart
  echo "restart stop"
	exit 1

elif [ $1 = "user" ]; then
  tmux kill-session -t userServer
  echo "userServer stop"
	exit 1

elif [ $1 = "sweep" ]; then
  tmux kill-session -t sweepServer
  echo "sweepServer stop"
	exit 1

elif [ $1 = "blockchain" ]; then
    for i in $(seq 1 2);
    do
       tmux kill-session -t blockChainIndexer_"$i"
    done
    echo "blockChainIndexer stop"
exit 1

elif [ $1 = "withdraw" ]; then
for i in $(seq 1 14);
do
  tmux kill-session -t withdrawObserver_"$i"
done
 echo "withdrawObserver stop"
exit 1

elif [ $1 = "deposit" ]; then
for i in $(seq 1 14);
do
  tmux kill-session -t depositObserver_"$i"
done
 echo "depositObserver stop"
exit 1

elif [ $1 = "approval" ]; then

    for i in $(seq 1 4);
    do
      tmux kill-session -t approvalObserver_"$i"
    done
 echo "approvalObserver stop"
  exit 1

elif [ $1 = "sign" ]; then
  tmux kill-session -t signServer
  echo "signServer stop"
	exit 1

elif [ $1 = "all" ]; then
  tmux kill-window -t configServer
	tmux kill-window -t authServer
	tmux kill-window -t userServer
	tmux kill-window -t signServer
  tmux kill-window -t sweepServer

  #kill-session
    for i in $(seq 1 14);
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
  echo "Stopped all server"
fi;
