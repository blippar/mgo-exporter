# Docker Compose Test Stack
## Testing
#### Emulate replicaLag

1. Lock Secondary
```sh
docker-compose exec secondary mongo --quiet --eval '
    printjson(db.fsyncLock())
'
```

2. Insert on Primary
```sh
docker-compose exec primary mongo --quiet --eval '
    printjson(db.collection.insert({item: "Beer", qty: 25}))
'
```

3. Unlock Secondary
```sh
docker-compose exec secondary mongo --quiet --eval '
    printjson(db.fsyncUnlock())
'
```

4. Print replicaSet secondary lag
```sh
docker-compose exec primary mongo --quiet --eval '
    rs.printSlaveReplicationInfo()
'
```
