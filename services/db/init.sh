#!/bin/sh
set -e

echo "waiting for postgres..."
until pg_isready -h "$PGHOST" -U "$PGUSER" -d "$PGDATABASE" > /dev/null 2>&1; do
  sleep 1
done
echo "postgres is ready"

echo "running migrations..."
psql -v ON_ERROR_STOP=1 -f /sql/migrations/000001_init.up.sql
echo "migrations done"

ROW_COUNT=$(psql -v ON_ERROR_STOP=1 -tAc "SELECT count(*) FROM providers")
if [ "$ROW_COUNT" = "0" ]; then
  echo "empty database, running seed..."
  psql -v ON_ERROR_STOP=1 -f /sql/seed.sql
  echo "seed done"
else
  echo "data already exists, skipping seed"
fi
