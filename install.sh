#!/bin/bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

if [[ -f "/usr/local/bin/gops" ]]; then
  echo "/usr/local/bin/gops already exists, removing"
  sudo rm /usr/local/bin/gops
fi

echo sudo ln -s "$SCRIPT_DIR/gops-stub.sh" "/usr/local/bin/gops"
sudo ln -s "$SCRIPT_DIR/gops-stub.sh" "/usr/local/bin/gops"
