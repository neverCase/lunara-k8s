#!/usr/bin/env bash

file="default.conf"
cat ${file}
rm ${file}
touch ${file}

if [[ "$ENV_REDIS_MASTER" ]] & [[ "$ENV_REDIS_MASTER_PORT" ]]; then
	echo "slaveof $ENV_REDIS_MASTER $ENV_REDIS_MASTER_PORT" >> ${file}
fi

if [[ "$ENV_REDIS_DIR" ]]
then
	echo "dir $ENV_REDIS_DIR" >> ${file}
fi

if [[ "$ENV_REDIS_DBFILENAME" ]]
then
	echo "dbfilename $ENV_REDIS_DBFILENAME" >> ${file}
fi

redis-server ${file}
