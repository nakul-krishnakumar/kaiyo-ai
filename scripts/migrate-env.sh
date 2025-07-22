#!/bin/bash
# This script sets up the environment for running the Goose database migration toolkit.
# Make sure you have Goose installed and available in your PATH.

echo "Setting up the environment for goose database migration"

echo "Enter database driver to use (default: postgres): "
read -r DB_DRIVER
if [ -z "$DB_DRIVER" ]; then
    DB_DRIVER="postgres"
fi
export GOOSE_DRIVER=$DB_DRIVER

echo "Enter database connection string (e.g., user:password@tcp(localhost:3306)/dbname): "
read -r DB_CONNECTION_STRING
export GOOSE_DBSTRING=$DB_CONNECTION_STRING

echo "Enter the path to the migration files (default: internal/database/migrations): "
read -r GOOSE_MIGRATION_DIR
if [ -z "$GOOSE_MIGRATION_DIR" ]; then
    GOOSE_MIGRATION_DIR="internal/database/migrations"
fi
# Ensure the migration directory exists
if [ ! -d "$GOOSE_MIGRATION_DIR" ]; then
    echo "Migration directory does not exist. Creating it at $GOOSE_MIGRATION_DIR"
    mkdir -p "$GOOSE_MIGRATION_DIR"
fi
# Export the migration directory
export GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR


echo ""
echo "Environment ready! Variables set:"
echo "GOOSE_DRIVER=$GOOSE_DRIVER"
echo "GOOSE_DBSTRING=$GOOSE_DBSTRING"
echo "GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR"
