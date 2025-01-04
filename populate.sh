#!/bin/bash

for shard in localhost:5000 localhost:5001; do
    echo $shard
    for i in {1..1000}; do
        curl "http://$shard/set?key=key-$RANDOM&value=value-$RANDOM"
    done
done