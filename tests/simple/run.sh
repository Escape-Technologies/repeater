#!/usr/bin/env bash

repo="$(dirname "$(dirname "$(dirname "${0}")")")"

(
    cd "${repo}" || exit 1
    go run "cmd/repeater/repeater.go"
)
