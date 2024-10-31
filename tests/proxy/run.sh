#!/usr/bin/env bash

export ESCAPE_REPEATER_PROXY_URL="http://localhost:9999"

repo="$(dirname "$(dirname "$(dirname "${0}")")")"

(
    cd "${repo}" || exit 1
    go run "cmd/repeater/repeater.go"
)
