#!/bin/bash

# Shell stub for development (place it into path and adjusts paths)
REPO="$HOME/git/mantro/goops"


set -euo pipefail

DIR="$(pwd)"

cd "$REPO"
go build

cd "$DIR"
"$REPO/goops" $@
