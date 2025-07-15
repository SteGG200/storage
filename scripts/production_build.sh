#!/usr/bin/env bash

set -e

cd frontend && pnpm install && pnpm build && cd ..

go build -o build/storage -tags production .