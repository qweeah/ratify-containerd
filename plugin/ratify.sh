#!/usr/bin/env bash

name=$2
digest=$4
export HOME=/root
if [ ! -f ~/.ratify/config.json ]; then
    echo "ratify config file not found. failing open"
    exit 0
fi

echo "=== input ==="
echo "1: $1"
echo "2: $2"
echo "3: $3"
echo "4: $4"
echo "==="

ratifyOutput=$(~/.ratify/bin/ratify verify -c ~/.ratify/config.json -s $name --digest $digest)
isSuccess=$(echo $ratifyOutput | jq .isSuccess)
echo $ratifyOutput | jq . 
if [ "$isSuccess" = "true" ]; then
    echo "ratify verification succeeded"
    exit 0
fi
exit 1