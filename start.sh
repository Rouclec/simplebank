#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database postgresql://postgres:%3FgR6%3A8Qm8e%21o4c%7C-99Y_I7%23V5%7D%3EX@simplebank.cz9maqdjakpb.us-east-1.rds.amazonaws.com:5432/simplebank -verbose up

echo "start the app"
exec "$@"