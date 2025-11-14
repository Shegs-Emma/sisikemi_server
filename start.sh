#!/bin/sh

echo "Running DB migrations..."
migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Starting Go app..."
exec /app/main
