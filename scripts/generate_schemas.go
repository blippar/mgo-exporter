package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/blippar/mgo-exporter/exporter"
	"github.com/blippar/mgo-exporter/mongo/model"
	"github.com/urakozz/go-json-schema-generator"
)

const schemaPath = "./documentation/schemas"

var jsonSchema = map[string]interface{}{
	"message.json":            &exporter.Message{},
	"mongo.json":              &exporter.ServerInfo{},
	"serverStatus.json":       &model.ServerStatus{},
	"dbStats.json":            &model.DBStats{},
	"replStatus.json":         &model.ReplStatus{},
	"replStatus.members.json": &model.ReplMember{},
	"nodeReplInfo.json":       &exporter.ReplicaSetInfo{},
}

func main() {

	for name, schema := range jsonSchema {
		sPath := path.Join(schemaPath, name)
		log.Println("Generating JSON-Schema:", sPath)
		s := generator.Generate(schema)
		err := ioutil.WriteFile(sPath, []byte(s), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
