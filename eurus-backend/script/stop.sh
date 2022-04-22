#!/bin/sh

if [ $# -eq 0 ] ; then
    echo 'Please input sh run.sh [ all | config | auth | user | withdraw | deposit | approval | restart | userObs | blockchain | sweep | merchantAdmin | publicData ]'
    exit 1
fi

RUN=0

killConfig() {
    tmux kill-session -t configServer
}

killAuth() {
    tmux kill-session -t authServer
}

killUser() {
    tmux kill-session -t userServer
}

killUserObserver() {
    tmux kill-session -t userObserver
}

killSweep() {
    tmux kill-session -t sweepServer
}

killWithdraw() {
    for i in $(seq 1 7);
    do
        tmux kill-session -t withdrawObserver_"$i"
    done
}

killRestart() {
   for i in $(seq 1 5);
    do
       tmux kill-window -t restart
    done

}

killApprove() {
    for i in $(seq 1 4);
    do
      tmux kill-session -t approvalObserver_"$i"
    done
}

killDeposit() {
    for i in $(seq 1 7);
    do
      tmux kill-session -t depositObserver_"$i"
    done
}

killBlockChain() {
    for i in $(seq 1 1);
    do
        tmux kill-session -t  blockChainIndexer_"$i"
    done
}


killMerchantAdmin() {
    tmux kill-session -t merchantAdminServer
}

killPublicData() {
    tmux kill-session -t publicData
}

if [ $1 = "config" ] || [ $1 = "all" ]; then
    RUN=1
    killConfig
fi

if [ $1 = "auth" ] || [ $1 = "all" ]; then
    RUN=1
    killAuth
fi

if [ $1 = "user" ] || [ $1 = "all" ]; then 
    RUN=1
    killUser
fi

if [ $1 = "withdraw" ] || [ $1 = "all" ]; then
    RUN=1
    killWithdraw
fi

if [ $1 = "approval" ] || [ $1 = "all" ]; then
    RUN=1
    killApprove
fi

if [ $1 = "deposit" ] || [ $1 = "all" ]; then 
    RUN=1
    killDeposit
fi

if [ $1 = "restart" ] || [ $1 = "all" ]; then
    RUN=1
    killRestart
fi

if [ $1 = "userObs" ] || [ $1 = "all" ]; then
    RUN=1
    killUserObserver
fi

if [ $1 = "sweep" ] || [ $1 = "all" ]; then
    RUN=1
    killSweep
fi

if [ $1 = "blockchain" ] || [ $1 = "all" ]; then
    RUN=1
    killBlockChain
fi

if [ $1 = "merchantAdmin" ] || [ $1 = "all" ]; then
    RUN=1
    killMerchantAdmin
fi

if [ $1 = "publicData" ] || [ $1 = "all" ]; then
    RUN=1
    killPublicData
fi

if [ $RUN -ne 1 ]; then
    echo "Invalid option"
    exit 1
else
    echo "Stopped"
fi

exit 0
