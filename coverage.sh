#!/usr/bin/env bash

set -e

# http://stackoverflow.com/a/21142256/2055281

echo "mode: atomic" > coverage.out
echo $@

for d in $@; do
    go test  -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        echo "$(pwd)"
        cat profile.out | grep -v "mode: " >> coverage.out
        rm profile.out
    fi
done
