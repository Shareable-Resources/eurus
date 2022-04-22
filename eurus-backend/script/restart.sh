# echo [a]uthServer --config ../config:./run.sh auth | awk '{split($0, a,":"); print a[2]}'

RUN_FOLDER=/home/ubuntu/build/src/eurus-backend/
while :
do
cd $RUN_FOLDER
#auth server
if [ $(ps aux | grep "/[a]uthServer" | wc -l) = "0" ]; then
    echo 'auth process restart'
    sh ./script/run.sh auth
else
    echo 'auth process running'
fi

#config server
if [ $(ps aux | grep "/[c]onfigServer" | wc -l) = "0" ]; then
    echo 'config process restart'
    sh ./script/run.sh config
else
    echo 'config process running'
fi

#user server
if [ $(ps aux | grep "/[u]serServer" | wc -l) = "0" ]; then
    sh ./script/run.sh user
else
    echo 'user process running'
fi

#approval server
for i in $(seq 1 4); do
  if [ $(ps aux | grep "[A]pprovalObserver_$i.json" | wc -l) = "0" ]; then
    sh ./script/run.sh approval $i
else
    echo 'approval process running'
fi
done

#withdraw server
for i in $(seq 1 7); do
  if [ $(ps aux | grep "[W]ithdrawObserverConfig_$i.json" | wc -l) = "0" ]; then
    sh ./script/run.sh withdraw $i
else
    echo 'approval process running'
fi
done

#deposit server
for i in $(seq 1 7); do
  if [ $(ps aux | grep "[D]epositObserverConfig_$i.json" | wc -l) = "0" ]; then
    sh ./script/run.sh deposit $i
else
    echo 'deposit process running'
fi
done

#blockchain indexer server
for i in $(seq 1 1); do
  if [ $(ps aux | grep "[B]lockChainIndexerConfig_$i.json" | wc -l) = "0" ]; then
    sh ./script/run.sh blockchain $i
else
    echo 'deposit process running'
fi
done

#sign server
if [ $(ps aux | grep "/[s]ignServer" | wc -l) = "0" ]; then
    sh ./script/run.sh sign
else
    echo 'signServer running'
fi

#userObserver
if [ $(ps aux | grep "/[u]serObserver" | wc -l) = "0" ]; then
    sh ./script/run.sh userObs
else
    echo 'userObserver running'
fi

#sweep server
if [ $(ps aux | grep "/[s]weepServer" | wc -l) = "0" ]; then
    sh ./script/run.sh sweep
else
    echo 'sweep process running'
fi

#sweep server
if [ $(ps aux | grep "/[k]ycServer" | wc -l) = "0" ]; then
    sh ./script/run.sh kyc
else
    echo 'kyc proccess running'
fi

if [ $(ps aux | grep "/[p]asswordServer" | wc -l) = "0" ]; then
    sh ./script/run.sh password
else
    echo 'passwordServer proccess running'
fi

sleep 5
done
