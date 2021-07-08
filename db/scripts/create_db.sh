#!/bin/bash

# This will fail if we have a db, else it will create one
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER postgres;
    CREATE DATABASE article;
    GRANT ALL PRIVILEGES ON DATABASE article TO postgres;
EOSQL
