#!/bin/sh
set -e

echo "run db migration"

# 假设你已经将 SecretString 存储在 SECRET_STRING 变量中
SECRET_STRING='{"DB_SOURCE":"postgresql://root:e72ZEDyYtruZfNzNO3r5@simple-bank.clgu0eqaq7ju.ap-southeast-2.rds.amazonaws.com:5432/simple_bank","DB_DRIVER":"postgres","SERVER_ADDRESS":"0.0.0.0:8080","ACCESS_TOKEN_DURATION":"15m","TOKEN_SYMETRIC_KEY":"beee0ad9c7bbd5ebb8bafb47b28ef325"}'

# 使用 jq 提取 DB_SOURCE
export DB_SOURCE=$(echo $SECRET_STRING | jq -r .DB_SOURCE)

# 检查 DB_SOURCE 是否被正确设置
if [ -z "$DB_SOURCE" ]; then
    echo "Error: DB_SOURCE is not set!"
    exit 1
fi

echo "DB_SOURCE: $DB_SOURCE"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
