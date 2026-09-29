#!/usr/bin/env bash

cd "$(dirname "$0")"

docker volume inspect vb-events_postgres_data &>/dev/null || docker volume create vb-events_postgres_data

docker rm -f vbevents-postgres 2>/dev/null

docker run -d \
  --name vbevents-postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=vbevents \
  -e POSTGRES_PASSWORD=vbevents \
  -e POSTGRES_DB=events \
  -v vb-events_postgres_data:/var/lib/postgresql/data \
  postgres:16-alpine
