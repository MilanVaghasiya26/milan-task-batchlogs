#!/bin/sh
set -e

create() {
  go run ../gogen/structure.go "$A1" "$A2" "$A3" "$A4" "$A5"
}

A1=$1
A2=$2
A3=$3
A4=$4
A5=$5

case "$1" in
'create')
  echo "Generate new structure!"
  create
  go fmt ./...
  echo "DONE"
  ;;
*)
  echo "Invalid input. Usage: $0 create" >&2
    exit 1
    ;;
esac
