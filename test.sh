#!/bin/bash
set -eu

echo Test Package
go test -v ./pkg/...

echo Test Examples
./test-examples.sh
