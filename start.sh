#!/bin/sh

set -e

# Start the first process
echo "Running db migration"
# shellcheck disable=SC2039
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the second process
echo "start the app"
exec "$@"