#!/usr/bin/env bash

defaultConf="redis.conf"
file="default.conf"
if [[ "$ENV_REDIS_CONF" ]]
then
	file="$ENV_REDIS_CONF"
else
    echo "error: no ENV_REDIS_CONF"
    exit
fi

cp ${defaultConf} ${ENV_REDIS_CONF}
rm ${file}
touch ${file}

echo "save 900 1" >> ${file}
echo "save 300 10" >> ${file}
echo "save 60 10000" >> ${file}

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

shutdownSave() {
   redis-cli shutdown save
}

trap "echo 'get the signal,redis-server would shut down and save before release container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

redis-server ${file} &

wait
