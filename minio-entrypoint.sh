#!/bin/sh
set -e

echo "Starting MinIO server..."
/usr/bin/docker-entrypoint.sh server /data --console-address ":9001" &
MINIO_PID=$!

echo "Waiting for MinIO to be ready..."
sleep 10

echo "Setting up MinIO client..."
mc alias set local http://localhost:9000 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"

echo "Creating panoramas bucket..."
mc mb local/panoramas --ignore-existing

echo "Setting bucket policy to public..."
mc anonymous set public local/panoramas

echo "MinIO setup completed - bucket panoramas is now public"

# Wait for MinIO to exit
wait $MINIO_PID