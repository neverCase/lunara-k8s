#!/usr/bin/env bash

redis-server --slaveof ${ENV_REDIS_MASTER} ${ENV_REDIS_MASTER_PORT}
