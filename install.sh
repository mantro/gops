#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [[ -f "/usr/local/bin/gops" ]]; then
  echo "/usr/local/bin/gops already exists"
  exit 1

fi

echo sudo ln -s "$SCRIPT_DIR/goops-stub.sh" "/usr/local/bin/gops"
sudo ln -s "$SCRIPT_DIR/goops-stub.sh" "/usr/local/bin/gops"

