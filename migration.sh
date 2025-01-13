#!/bin/sh
set -e

migrate() {
  go run data_model/migration.go "$A1" "$A2" "$A3"
}

A1=$1
A2=$2
A3=$3

case "$1" in
'migration')
  echo "Generate migration file!"
  migrate
  ;;
'migrate')
  echo "Migrate database version"
  migrate
  ;;
'seeder')
  echo "Database seeding"
  migrate
  ;;
*)
  BadInput
  ;;
esac
