#!/bin/sh
set -e

echo "run db migration"

source /app/app.env
cat /app/app.env

echo "DB_SOURCE: $DB_SOURCE"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
