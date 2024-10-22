#!/bin/sh
set -e

echo "run db migration"
source /app/app.env

if [ ! -f /app/app.env ]; then
    echo "Error: app.env file not found!"
    exit 1
fi

echo /app/app.env


# 打印 DB_SOURCE 以调试
echo "DB_SOURCE: $DB_SOURCE"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
