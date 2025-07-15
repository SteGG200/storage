#!/usr/bin/env bash

FILE=${FILE:-data.db}

go run ./db/migrations/migrate.go -database ${FILE}