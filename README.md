# mgo-exporter

A simple MongoDB stats exporter with support for ReplicaSet cluster.

### Exported command output

Command             | Exported as
--------------------|---------------------
`db.serverStatus()` | `serverStatus`
`rs.Status()`       | `replStatus`
`db.stats()`        | `dbStats`

### Exporter tested on
- [x] MongoDB 3.0 with Logstash 5.4
