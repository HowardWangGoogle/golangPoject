#!/bin/sh
set -e

echo "run db migration"

source /app/.env.docker
cat /app/.env.docker

echo "DB_SOURCE: $DB_SOURCE"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"



