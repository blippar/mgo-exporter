# Emulate replicaLag

- Lock Secondary
```sh
docker-compose exec secondary mongo --quiet --eval 'printjson(db.fsyncLock())'
```

- Insert on Primary
```sh
docker-compose exec primary mongo --quiet --eval 'printjson(db.collection.insert({item: "Beer", qty: 25}))'
```

- Unlock Secondary
```sh
docker-compose exec secondary mongo --quiet --eval 'printjson(db.fsyncUnlock())'
```

- Print replicaSet secondary lag
```sh
docker-compose exec primary mongo --quiet --eval 'rs.printSlaveReplicationInfo()'
```
