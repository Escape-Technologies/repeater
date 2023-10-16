#!/bin/bash

IMG=bufbuild/buf:1.27.0
DIR="/repo"

cd "$(dirname "$BASH_SOURCE[0]")"

docker run -v "$(pwd)/..:${DIR}" \
    --rm \
    --user "$(id -g):$(id -g)" \
    -e HOME="/tmp/" \
    --workdir="${DIR}" \
    --entrypoint=/bin/sh \
    "${IMG}" \
    -c '
cd protocol
buf mod update .
buf format -w
buf generate --template buf.gen.yaml --path repeater

# Recursively remove all occurences of ,omitempty
find .. -type f -name "*.go" -exec sed -i 's/,omitempty//g' {} \;
'
