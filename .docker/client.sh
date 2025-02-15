#!/usr/bin/env sh

set -eu

dbhost="${1:?missing dbhost}"
dbname='pgstore_test'
dbuser='pgstore'

# Wait for db server to be ready, with some limits.

num_attempts=0
until PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${dbhost}" -U "${dbuser}" "${dbname}" -c '\q'; do
	num_attempts=$((num_attempts+1))
	if [ $num_attempts -gt 10 ]; then
		>&2 echo "ERROR: max attempts exceeded"
		exit 1
	fi

	>&2 echo "db is unavailable now, sleeping"
	sleep 1
done
>&2 echo "db is up"

echo "testing against live DB"
make check
