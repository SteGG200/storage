#!/usr/bin/env bash

set -e 

FILE=${FILE:-data.db}

if test -f "${FILE}"; then
	echo "Backing up as ${FILE}.bak"
	mv ${FILE} "${FILE}.bak"
fi

touch ${FILE}

go run "$(pwd)/db/migrations/migrate.go" -database ${FILE}