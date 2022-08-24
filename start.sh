#!/bin/sh

set -e

echo "start the app"
exec "$@"

CompileDaemon --build="go build -o main main.go"  --command=./main