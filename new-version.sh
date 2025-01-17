#!/bin/bash

old_version=$(cat helm/Chart.yaml | grep 'version:' | awk '{print $2}')
SEMVER_REGEX="^([0-9]+)\.([0-9]+)\.([0-9]+)$"

if [[ "$old_version" =~ $SEMVER_REGEX ]]; then
    _major=${BASH_REMATCH[1]}
    _minor=${BASH_REMATCH[2]}
    _patch=${BASH_REMATCH[3]}
else
    echo "Invalid version in helm/Chart.yaml: $old_version"
    exit 1
fi

case $1 in
  patch)
    _patch=$(($_patch + 1))
    ;;
  minor)
    _minor=$(($_minor + 1))
    _patch=0
    ;;
  *)
    echo "Usage: $0 <minor|patch>"
    exit 1
    ;;
esac

new_version="${_major}.${_minor}.${_patch}"
echo "Bump done: $old_version -> $new_version"

cat <<EOF > helm/Chart.yaml
---
apiVersion: v2
name: escape-repeater
description: Escape repeater
type: application
version: ${new_version}
appVersion: ${new_version}
EOF

git add helm/Chart.yaml
git commit -m "v${new_version}"
git tag -a "v${new_version}" -m "v${new_version}"
git push
git push --tags

echo "Done !"
