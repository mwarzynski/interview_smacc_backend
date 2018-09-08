#!/bin/bash
set -e

export GATEWAY=`route -n | awk 'FNR==3{print $2}'`
export METRICS_HOST=$GATEWAY

if [ "$1" = 'run' ]; then
    exec /server
fi

exec "$@"
