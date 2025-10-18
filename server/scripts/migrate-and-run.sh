#!/usr/bin/env sh
set -eu

# migrate-and-run.sh
# Run goose migrations (with retry) then exec the provided command (or air by default)

MIGRATIONS_DIR=${MIGRATIONS_DIR:-/app/internal/database/migrations}
DB_URL=${DATABASE_URL:-}
MAX_ATTEMPTS=${MAX_ATTEMPTS:-10}
SLEEP_SECONDS=${SLEEP_SECONDS:-2}

run_migrations() {
  if [ -z "$DB_URL" ]; then
    echo "DATABASE_URL not set, skipping migrations"
    return
  fi

  echo "Running migrations from $MIGRATIONS_DIR against $DB_URL"
  attempt=1
  while [ $attempt -le $MAX_ATTEMPTS ]; do
    if /usr/local/bin/goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up; then
      echo "Migrations applied successfully"
      return 0
    fi
    echo "Migration attempt $attempt failed — retrying in ${SLEEP_SECONDS}s..."
    attempt=$((attempt + 1))
    sleep $SLEEP_SECONDS
  done

  echo "Migrations failed after $MAX_ATTEMPTS attempts"
  return 1
}

# Try to run migrations but don't fail startup hard in dev; exit code preserved
if run_migrations; then
  echo "Migrations finished"
else
  echo "Continuing startup despite migration failures"
fi

# If no args passed, default to air -- --port 8081
if [ "$#" -eq 0 ]; then
  set -- air -- --port 8081
fi

exec "$@"
