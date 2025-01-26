#!/usr/bin/env sh

FILE=${FILE:-data.db}

go run ./db/migrations/migrate.go -database ${FILE}