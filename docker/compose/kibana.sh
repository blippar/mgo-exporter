#!/bin/sh
set -e

echo "> Waiting for elastic to become available"
until curl -sS --fail "http://elastic:9200/_cat/health" >/dev/null 2>/dev/null; do
    sleep 1
done

echo "> Create kibana index"

curl -sS -XPUT "http://elastic:9200/.kibana/?pretty" -d '{
  "aliases": {},
  "mappings": {
    "config": {
      "properties": {
        "buildNum": {
          "type": "keyword"
        }
      }
    },
    "index-pattern": {
      "properties": {
        "fieldFormatMap": {
          "type": "text"
        },
        "fields": {
          "type": "text"
        },
        "intervalName": {
          "type": "keyword"
        },
        "notExpandable": {
          "type": "boolean"
        },
        "sourceFilters": {
          "type": "text"
        },
        "timeFieldName": {
          "type": "keyword"
        },
        "title": {
          "type": "text"
        }
      }
    },
    "visualization": {
      "properties": {
        "description": {
          "type": "text"
        },
        "kibanaSavedObjectMeta": {
          "properties": {
            "searchSourceJSON": {
              "type": "text"
            }
          }
        },
        "savedSearchId": {
          "type": "keyword"
        },
        "title": {
          "type": "text"
        },
        "uiStateJSON": {
          "type": "text"
        },
        "version": {
          "type": "integer"
        },
        "visState": {
          "type": "text"
        }
      }
    },
    "search": {
      "properties": {
        "columns": {
          "type": "keyword"
        },
        "description": {
          "type": "text"
        },
        "hits": {
          "type": "integer"
        },
        "kibanaSavedObjectMeta": {
          "properties": {
            "searchSourceJSON": {
              "type": "text"
            }
          }
        },
        "sort": {
          "type": "keyword"
        },
        "title": {
          "type": "text"
        },
        "version": {
          "type": "integer"
        }
      }
    },
    "dashboard": {
      "properties": {
        "description": {
          "type": "text"
        },
        "hits": {
          "type": "integer"
        },
        "kibanaSavedObjectMeta": {
          "properties": {
            "searchSourceJSON": {
              "type": "text"
            }
          }
        },
        "optionsJSON": {
          "type": "text"
        },
        "panelsJSON": {
          "type": "text"
        },
        "refreshInterval": {
          "properties": {
            "display": {
              "type": "keyword"
            },
            "pause": {
              "type": "boolean"
            },
            "section": {
              "type": "integer"
            },
            "value": {
              "type": "integer"
            }
          }
        },
        "timeFrom": {
          "type": "keyword"
        },
        "timeRestore": {
          "type": "boolean"
        },
        "timeTo": {
          "type": "keyword"
        },
        "title": {
          "type": "text"
        },
        "uiStateJSON": {
          "type": "text"
        },
        "version": {
          "type": "integer"
        }
      }
    },
    "url": {
      "properties": {
        "accessCount": {
          "type": "long"
        },
        "accessDate": {
          "type": "date"
        },
        "createDate": {
          "type": "date"
        },
        "url": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 2048
            }
          }
        }
      }
    },
    "server": {
      "properties": {
        "uuid": {
          "type": "keyword"
        }
      }
    },
    "timelion-sheet": {
      "properties": {
        "description": {
          "type": "text"
        },
        "hits": {
          "type": "integer"
        },
        "kibanaSavedObjectMeta": {
          "properties": {
            "searchSourceJSON": {
              "type": "text"
            }
          }
        },
        "timelion_chart_height": {
          "type": "integer"
        },
        "timelion_columns": {
          "type": "integer"
        },
        "timelion_interval": {
          "type": "keyword"
        },
        "timelion_other_interval": {
          "type": "keyword"
        },
        "timelion_rows": {
          "type": "integer"
        },
        "timelion_sheet": {
          "type": "text"
        },
        "title": {
          "type": "text"
        },
        "version": {
          "type": "integer"
        }
      }
    }
  },
  "settings": {
    "index": {
      "number_of_replicas": "1",
      "number_of_shards": "1"
    }
  }
}'

cd /opt/kibana
/opt/kibana/import.sh
