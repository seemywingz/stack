#!/bin/bash
ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Move stack Binary to /usr/local/bin
cp ${ROOT_DIR/scripts/bin}/stack /usr/local/bin
