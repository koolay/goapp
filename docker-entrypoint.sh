#!/bin/bash
set -e
cd $APP_HOME

if [ "$1" = 'goapp' -a "$(id -u)" = '0' ]; then
	echo "[entrypoint] gosu run app"
	exec gosu app "$@" -c $APP_HOME
fi
exec "$@"

