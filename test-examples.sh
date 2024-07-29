#!/bin/bash
set -eu
for d in examples/*/; do
  cd "$d"
  echo
  echo "pwd: $(pwd)"
  go test -v ./...
  cd ../..
done
