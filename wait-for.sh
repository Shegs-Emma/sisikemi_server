#!/bin/sh

# wait-for.sh host:port -- command
host="$1"
shift
until nc -z ${host%:*} ${host#*:}; do
  echo "Waiting for $host..."
  sleep 1
done

exec "$@"
