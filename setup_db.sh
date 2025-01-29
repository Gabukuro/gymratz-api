#!/bin/bash

psql service=dev -c "SELECT 1 FROM pg_database WHERE datname = 'gymratz-api'" | grep -q 1 > /dev/null
[ $? -eq 0 ] && DB_EXISTS=1 || DB_EXISTS=0

[ $DB_EXISTS -eq 0 ] && psql service=dev -c 'CREATE DATABASE "gymratz-api"'

sql-migrate up

[ $DB_EXISTS -eq 0 ] && sql-migrate up -env='seeder' || exit 0
