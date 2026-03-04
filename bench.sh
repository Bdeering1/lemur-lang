#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 [directory]"
  exit 1
fi

DIR=$1

echo "Priming CPU..."
go test "$DIR" -bench=. -count=2 > /dev/null

echo
go test "$DIR" -bench=. -count=10 | tee /dev/stderr | benchstat /dev/stdin
