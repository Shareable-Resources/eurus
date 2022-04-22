#!/bin/bash

SCRIPT_NAME=$(basename $0)
if [ $# -lt 1 ]; then
	echo "Usage: $SCRIPT_NAME [config file path]"
	exit 1
fi


source $1

while true; do
	CUR_TIME=$(date +%s%N)
	psql postgresql://${KEEP_ALIVE_USER}:${KEEP_ALIVE_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres -c "UPDATE keep_alive_config SET last_modified_time = $CUR_TIME WHERE id = 1;"
	sleep 0.3
done