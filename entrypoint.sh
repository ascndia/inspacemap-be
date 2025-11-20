#!/bin/sh
set -e 

echo "🚀 [Entrypoint] Starting deployment process..."

echo "🌱 [Entrypoint] Running Database Seeder..."
./seeder

echo "🔥 [Entrypoint] Starting Go Backend Server..."
exec ./main