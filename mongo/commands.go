package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/blippar/mgo-exporter/mongo/model"
)

const adminDB = "admin"

// GetServerStatus ...
func GetServerStatus(session *mgo.Session) (*model.ServerStatus, error) {

	result := &model.ServerStatus{}

	err := session.DB(adminDB).Run(
		bson.D{
			{Name: "serverStatus", Value: 1},
			{Name: "locks", Value: 0},
			{Name: "repl", Value: 0},
			{Name: "metrics", Value: 0},
			{Name: "cursors", Value: 0},
			{Name: "storageEngine", Value: 0},
			{Name: "dur", Value: 0},
			{Name: "extra_info", Value: 0},
			{Name: "backgroundFlushing", Value: 0},
		}, result)
	if err != nil {
		result.OK = 0
		result.Error = err.Error()
	}
	return result, err
}

// GetDBStats ...
func GetDBStats(session *mgo.Session, db string) (*model.DBStats, error) {

	result := &model.DBStats{}

	err := session.DB(db).Run(
		bson.D{
			{Name: "dbStats", Value: 1},
			// {Name: "scale", Value: 1024}, // scale size to kilobytes rather than bytes
		}, result)
	if err != nil {
		result.DB = db
		result.OK = 0
		result.Error = err.Error()
	}
	return result, err
}

// GetReplStatus ...
func GetReplStatus(session *mgo.Session) (*model.ReplStatus, error) {

	result := &model.ReplStatus{}

	err := session.DB(adminDB).Run(
		bson.D{
			{Name: "replSetGetStatus", Value: 1},
		}, result)

	if err != nil {
		result.OK = 0
		result.Error = err.Error()
		return result, err
	}

	result.JSONMembers = make(map[int]*model.ReplMember)
	for _, v := range result.Members {
		result.JSONMembers[v.ID] = v
	}

	return result, nil
}
