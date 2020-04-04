#!/usr/bin/env bash

defaultConf="/redis.conf"
if [[ -z "$ENV_REDIS_CONF" ]]
then
    echo "error: no env ENV_REDIS_CONF"
    exit
fi

cp ${defaultConf} ${ENV_REDIS_CONF}

if [[ "$ENV_REDIS_MASTER" ]] && [[ "$ENV_REDIS_MASTER_PORT" ]]; then
	sed -i "s/# replicaof <masterip> <masterport>/replicaof ${ENV_REDIS_MASTER} ${ENV_REDIS_MASTER_PORT}/g"  ${ENV_REDIS_CONF}
fi

if [[ "$ENV_REDIS_DIR" ]]
then
	sed -i "s#dir ./#dir ${ENV_REDIS_DIR}#g"  ${ENV_REDIS_CONF}
fi

if [[ "$ENV_REDIS_DBFILENAME" ]]
then
	sed -i "s/dbfilename dump.rdb/dbfilename ${ENV_REDIS_DBFILENAME}/g" ${ENV_REDIS_CONF}
fi

shutdownSave() {
   redis-cli shutdown save
}

trap "echo 'get the signal,redis-server would shut down and save before release container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

redis-server ${ENV_REDIS_CONF} &

wait