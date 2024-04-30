#!/bin/sh
set -e

if [[ -z "$DATABASE_URL" ]]; then
    echo "Must provide DATABASE_URL in environment" 1>&2
    exit 1
fi

/app/migrate -source file:///app/migrations/ -database $DATABASE_URL up

/app/xbankapi
