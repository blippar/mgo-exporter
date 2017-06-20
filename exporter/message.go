package exporter

import (
	"time"

	"github.com/blippar/mgo-exporter/mongo/model"
)

// Message ...
type Message struct {
	Time         time.Time           `json:"time"`
	Mongo        *ServerInfo         `json:"mongo"`
	ServerStatus *model.ServerStatus `json:"serverStatus,omitempty"`
	DBStats      *model.DBStats      `json:"dbStats,omitempty"`
	Repl         *model.ReplStatus   `json:"replStatus,omitempty"`
	NodeReplInfo *ReplicaSetInfo     `json:"nodeReplInfo,omitempty"`
	Type         string              `json:"type"`
}

// ServerInfo ...
type ServerInfo struct {
	Host       string `json:"host"`
	ReplicaSet string `json:"replSet,omitempty"`
	Connected  bool   `json:"connected"`
	Error      string `json:"error,omitempty"`
}

// ReplicaSetInfo ...
type ReplicaSetInfo struct {
	Set           string        `json:"set"`
	Name          string        `json:"name"`
	State         int           `json:"state"`
	StateStr      string        `json:"stateStr"`
	OPTime        int64         `json:"optime"`
	OPTimeDate    *time.Time    `json:"optimeDate"`
	Uptime        int           `json:"uptime"`
	ConfigVersion int           `json:"configVersion"`
	OPTimeDiff    map[int]int64 `json:"optimeDiff"`
	OPTimeLag     int64         `json:"optimeLag"`
}
