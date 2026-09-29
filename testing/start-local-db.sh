#!/usr/bin/env bash

cd "$(dirname "$0")"

docker volume inspect postgres_data &>/dev/null || docker volume create postgres_data

docker rm -f vbevents-postgres 2>/dev/null

docker run -d \
  --name vbevents-postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=vb-events \
  -e POSTGRES_PASSWORD=vbevents \
  -e POSTGRES_DB=events \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:16-alpine
