#!/bin/sh
set -e

echo "run db migration"
source "/app/app.env"

DB_SOURCE=postgresql://root:e72ZEDyYtruZfNzNO3r5@simple-bank.clgu0eqaq7ju.ap-southeast-2.rds.amazonaws.com:5432/simple_bank
DB_DRIVER=postgres
SERVER_ADDRESS=0.0.0.0:8080
ACCESS_TOKEN_DURATION=15m
TOKEN_SYMETRIC_KEY=beee0ad9c7bbd5ebb8bafb47b28ef325


# 打印 DB_SOURCE 以调试
echo "DB_SOURCE: $DB_SOURCE"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
