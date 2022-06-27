#!/bin/sh
OPTIONS=`vagrant ssh-config | grep -v '^Host ' | awk -v ORS=' ' '{print "-o " $1 "=" $2}'`

scp ${OPTIONS} "$@" || echo "Transfer failed. Did you use 'default:' as the target?"

