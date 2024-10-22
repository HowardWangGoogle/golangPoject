#!/bin/sh
set -e

echo "run db migration"

# 加载环境变量
if [ ! -f /app/app.env ]; then
    echo "Error: app.env file not found!"
    exit 1
fi
source /app/app.env

# 创建 JSON 字符串
env_json=$(jq -n \
    --arg db_source "$DB_SOURCE" \
    --arg db_driver "$DB_DRIVER" \
    --arg server_address "$SERVER_ADDRESS" \
    --arg access_token_duration "$ACCESS_TOKEN_DURATION" \
    --arg token_symetric_key "$TOKEN_SYMETRIC_KEY" \
    '{
        DB_SOURCE: $db_source,
        DB_DRIVER: $db_driver,
        SERVER_ADDRESS: $server_address,
        ACCESS_TOKEN_DURATION: $access_token_duration,
        TOKEN_SYMETRIC_KEY: $token_symetric_key
    }')

# 输出 JSON
echo "Extracted Environment Variables in JSON format:"
echo "$env_json"

# 运行迁移
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
