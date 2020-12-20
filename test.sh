#!/usr/bin/env bash
set -e
set -o pipefail

if ! curl -f localhost:8081/artifactory/api/system/ping &> /dev/null; then
  docker-compose up -d
fi

until curl -f localhost:8081/artifactory/api/system/ping &> /dev/null; do
  sleep 5
done
go test ./tests