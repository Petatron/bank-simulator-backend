#!/bin/sh

set -e

# Start the first process
echo "Running db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Start the second process
echo "start the app"
exec "$@"