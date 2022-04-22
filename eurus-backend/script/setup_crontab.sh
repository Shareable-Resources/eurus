#!/bin/sh
RUN_FOLDER=~/eurus/bin
TEST=~/eurus/bin/
TEMP=~/eurus/tmp
CONFIG_FOLDER=~/eurus/config
ERRORLOG_FOLDER=~/eurus/errorLog
TIMESTAMP=`date "+%Y-%m-%d"`
USER_NAME=ubuntu
#Local
#USER_NAME=Vince

mkdir -p "$TEMP"

echo "*/60 * * * * sudo $RUN_FOLDER/background --config $CONFIG_FOLDER/BackgroundServerConfig.json" > "$TEMP/exchangeRateCron"
echo "00 01 * * * $RUN_FOLDER/walletBackground --config $CONFIG_FOLDER/WalletBackgroundIndexerConfig.json >> $ERRORLOG_FOLDER/WalletBackgroundIndexer.log" > "$TEMP/walletBackgroundCron" 
echo "*/5 * * * * sudo $RUN_FOLDER/blockCypherBackground --config $CONFIG_FOLDER/BlockCypherBackgroundConfig.json" > "$TEMP/blockCypherCron"


cat "$TEMP/exchangeRateCron" "$TEMP/walletBackgroundCron" "$TEMP/blockCypherCron" > "$TEMP/allCronJob"
sudo crontab -u $USER_NAME $TEMP/allCronJob
echo "finished setting up crontab"
