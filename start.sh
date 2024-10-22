#!/bin/sh

set -e


echo "run db migration"

source /app/app.env
cat /app/app.env

# 解析 JSON 并导出环境变量
eval $(jq -r 'to_entries | .[] | "export \(.key)=\(.value)"' /app/app.env)

# 输出 DB_SOURCE 的值
echo "DB_SOURCE: $DB_SOURCE"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"