<p align="center">
<img src="documentation/banner.png"></img>
<h4 align="center">A simple MongoDB stats exporter with replicaSet support</h3>
</p>

## Features

- **Export server metrics** exposed by MongoDB
- **Export replicaSet status** exposed by MongoDB
- **Get insight from your replicaSet** with calculation of the largest lag between 2 members
- **Get information about the size of your databases**

## Install

Mongo exporter is available in many form for the Linux platform, you can install it via:

- **A system package** (deb & rpm) available for download on the [Latest Release Page](https://github.com/blippar/mgo-exporter/releases/latest)
- **A Docker Container** available on the [DockerHub](https://hub.docker.com/r/blippar/mgo-exporter/)
- **go get** by running `go get github.com/blippar/mgo-exporter/cmd/mgo-exporter` (requires Go 1.6 or newer)

## Quick Start

To see an example of the data being exported by the `mgo-exporter` once installed and assuming your MongoDB is accessible via `127.0.0.1:27017`, run:
```
mgo-exporter 127.0.0.1:27017
```

### Running with Docker

Running the exporter via Docker is as easy as running it in the command line, container takes the exact same options and environment variable as the binary version.

```
docker run --restart=always \
           --link mongo:mongo \
           --link logstash:logstash \
           -e MGOEXPORT_REPL="myReplicaSet" \
           -e MGOEXPORT_FORWARDER="logstash://logstash:2000" \
           blippar/mgo-exporter mongo:27017
```

### Running with a system package

By installing `mgo-exporter` via a system package, it will automatically add a `systemd` unit file as well as a `sysconfig` style configuration file. Once installed via one of the following command:
```sh
# Centos
sudo yum install mgo-exporter_*.x86_64.rpm
# Debian
sudo dpkg -i mgo-exporter_*.amd64.dev
```

You will be able to configuration your software by editing `/etc/sysconfig/mgo-exporter` which contains the following:
```sh
MGOEXPORT_MONGODB="localhost:27017"
MGOEXPORT_DATABASES="test local"
MGOEXPORT_REPL=""
MGOEXPORT_FORWARDER="file:///var/log/mgo-exporter-metrics.log"
MGOEXPORT_LOGFILE=/var/log/mgo-exporter.log
```

Once edited and configured you can run the `mgo-exporter` as any other daemon for systemd:
```sh
# Start the exporter
sudo systemctl start mgo-exporter
# Automatically start with the system
sudo systemctl enable mgo-exporter
```

## Configuration
#### Usage
```
Usage: mgo-exporter [--repl REPL] [--forwarder FORWARDER] [--logfile LOGFILE] [--verbose] [--quiet] MONGODB [DATABASE [DATABASE ...]]

Positional arguments:
  MONGODB                Mongo URI for the node to connect to
  DATABASE               database name to monitor

Options:
  --repl REPL, -r REPL   replicaSet name to monitor [env: MGOEXPORT_REPL]
  --forwarder FORWARDER, -f FORWARDER
                         forwarder URI to send messages to [env: MGOEXPORT_FORWARDER] [default: file:///dev/stdout?pretty]
  --logfile LOGFILE, -l LOGFILE
                         file to output logs to [env: MGOEXPORT_LOGFILE] [default: /dev/stderr]
  --verbose, -v          enable a more verbose logging
  --quiet, -q            enable quieter logging
  --help, -h             display this help and exit
```

#### Extends what is monitored

By default `mgo-exporter` will export most of the fields available by running `db.serverStatus()` on your MongoDB instance, though you can extend what it will export by asking it to also output information about your databases as well as your replicaSet.

For example the following command will export information about the database `test` as well as information about the `myReplSet` replicaSet.

```
mgo-exporter --replicaSet myReplSet 127.0.0.1:27017 test
```

To see exactly what is exposed at any point, check the [Specifications](#specifications) section of this documentation.

#### Forward your data

When no forwarder is explicitly specified, `mgo-exporter` will print the gathered information to `/dev/stdout`, in order for those data to be useful you might want to store them or send them to another entity to be processed.

##### Store data in a file

Right out of the box, you can write those data to a file (here `/var/log/mongo-metrics.log`) by running the exporter in the following way:
```
mgo-exporter --forwarder file:///var/log/mongo-metrics.log 127.0.0.1:27017
```

##### Export them to an ELK Stack

Another solution is to forward those data to a Logstash instance in order to store them in a Database, this repository contains sample configuration to use the exporter with an ELK (ElasticSearch, Logstash & Kibana) pipeline.

In this example we will use a Logstash instance setup using the configuration available [here](docker/logstash/config.conf) and this [mapping](docker/logstash/mapping/mgo-exporter.template.json), this configuration get the raw message exported by the `mgo-exporter` and split all `dbStats` elements in its own document (for a better indexing on ElasticSearch).

In order to forward your data to it, you will just need to run the following:
```
mgo-exporter --forwarder logstash://127.0.0.1:2000 127.0.0.1:27017
```

This is the setup we recommend and thus we have created a sample Kibana (5.4+) configuration for such setup so we can visualize those data cleanly. We created different visualizations as well as a Dashboard that can be imported using the [`kibana/import.sh`](docker/kibana/import.sh) script we've created. This script will create everything you need prefixed with the `mgo-` prefix for easy cleanup.


##### More forwarders

To see all available forwarder as well as how to configure them, check the [Available forwarder](#available-forwarders) section of this documentation.

## Running in production

The `mgo-exporter` has been created in order to monitor our production MongoDB and thus should be fitting if you are looking for something similar.

Though we recommend the following configuration for such a setup:
- Run the exporter alongside your MongoDB on the same host for better latency between the exporter and MongoDB
- Either use the Docker container or one of our system package while installing the exporter (our packages contains everything needed to supervise the exporter via `systemd`)
- Forward your data to an external source (we initially built this software to export MongoDB metrics to an ELK stack)

## Development
### Building it
All elements of `mgo-exporter` can be built via a `Makefile` at the root of this repository, available rules are the followings:
- `make exporter`: create a binary for your current system under `bin/mgo-exporter`
- `make static`: create a static binary targeted at Linux based system under `bin/mgo-exporter`
- `make rpm`: create a RPM package for `x86_64` system
- `make deb`: create a DEB package for `x86_64` system
- `make docker`: create the `blippar/mgo-exporter` Docker Container (tagged with the current version)
- `make dist`: create a RPM and a DEB package as well as the Docker Container
- `make generate_schemas`: regenerate JSON-Schema files in `documentations/schemas/*.json`
- `make clean`: remove all files that would be created by `make exporter rpm deb`

*NOTE*: system package created from this `Makefile` comes packages with a static binary of `mgo-exporter`, a `systemd` init script as well as `sysconfig` config file.

### Testing it
A complete `docker-compose` stack is available in this repository to test how the software would work, by running `docker-compose up -d`, this stack will create:

- 3 MongoDB setup in a replicaSet (named: `mgoExporterSet`) with 3 different priority
- 3 `mgo-exporter` setup to export each nodes information to a Logstash pipeline
- 1 Logstash setup to split DBStats elements in their own document and forward those documents to ElasticSearch
- 1 ElasticSearch to store the data forwarded by logstash
- 1 Kibana to show those data using our visualizations and dashboards
- 2 provisioning scripts:
    - One to setup the MongoDB replicaSet
    - The other to import index-pattern, visualizations and dashboard to Kibana

From there you can test how this stack would work in a `replicaSet` configuration by inserting / updating / deleting documents as well as emulate replicaLag. Some manual test used to create the visualizations currently present in this repository are available in the [`docker/compose/README.md`](docker/compose/README.md).

## Specifications
### Exported top level fields output

Exported fields | Command               | Format      | When
----------------|:---------------------:|:-----------:|--------------------------------------------------
`time`          | n/a                   | `date-time` | Always
`mongo`         | n/a                   | [schema][1] | Always
`serverStatus`  | [`db.serverStatus()`] | [schema][2] | When connected to `MONGODB`
`[]dbStats`     | [`db.stats()`]        | [schema][3] | When at least one `DATABASE` argument is passed and  
_               | _                     | _           | - with `--repl REPL`: when connected to the primary
_               | _                     | _           | - without `--repl`: when connected to `MONGODB`
`replStatus`    | [`rs.Status()`]       | [schema][4] | When connected and run with `--repl REPL`
`nodeReplInfo`  | generated             | [schema][5] | When connected and run with `--repl REPL`

[`db.serverStatus()`]: https://docs.mongodb.com/v3.0/reference/command/serverStatus/#dbcmd.serverStatus
[`db.stats()`]: https://docs.mongodb.com/v3.0/reference/command/dbStats/#dbcmd.dbStats
[`rs.Status()`]: https://docs.mongodb.com/manual/reference/command/replSetGetStatus/#dbcmd.replSetGetStatus

[1]: documentation/schemas/mongo.json
[2]: documentation/schemas/serverStatus.json
[3]: documentation/schemas/dbStats.json
[4]: documentation/schemas/replStatus.json
[5]: documentation/schemas/nodeReplInfo.json

### Available forwarders

Exporter   | Configuration                                | Description
-----------|----------------------------------------------|----------------------------------------------------------------------
`logstash` | `logstash://{logstash_addr}:{logstash_port}` | Send exported data to a Logstash TCP endpoint (compatible with `codec=>json`)
`file`     | `file://{file_path}[?pretty]`                | Store exported data in a logfile defined by `{file_path}`
_          | _                                            | If `?pretty` is specified, the `json` output will be prettified
