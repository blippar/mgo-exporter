#!/usr/bin/env bash
set -e

echo "> Waiting for mongo hosts to become available"
until mongo --quiet --host primary --eval '{ ping:1 }' > /dev/null; do
    sleep 1
done
until mongo --quiet --host fallback --eval "{ ping:1 }" > /dev/null; do
    sleep 1
done
until mongo --quiet --host secondary --eval "{ ping:1 }" > /dev/null; do
    sleep 1
done

echo "> Checking current replicaSet configuration"
not_init=$(mongo --host primary --quiet --eval 'printjson(rs.status())' | sed -nE 's/^[\t ]*"code" ?: ([0-9]*),?$/\1/p')

if [ "$not_init" != "" ] && [ "$not_init" -eq 94 ]; then
    echo -n "> Initilializing ReplicaSet: "
    mongo --quiet --host primary --eval 'printjson(rs.initiate(
       {
          _id: "mgoExporterSet",
          version: 1,
          members: [
             { _id: 0, host : "primary:27017", priority: 2 },
             { _id: 1, host : "fallback:27017", priority: 1 },
             { _id: 2, host : "secondary:27017", priority: 0 }
          ]
       }
    ))'
else
    echo "> ReplicaSet already initiliazed"
    exit 0
fi
