#!/bin/sh 

for mod_dir in $(find . -name go.mod | xargs -n 1 dirname | sort)
do
    cd $mod_dir
    echo $mod_dir
    go vet
    cd ../
done
