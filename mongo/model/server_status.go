package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ServerStatus ...
type ServerStatus struct {
	Host           string       `bson:"host,omitempty" json:"host,omitempty"`
	LocalTime      *time.Time   `bson:"localTime,omitempty" json:"localTime,omitempty"`
	PID            int          `bson:"pid,omitempty" json:"pid,omitempty"`
	Process        string       `bson:"process,omitempty" json:"process,omitempty"`
	Uptime         int          `bson:"uptime,omitempty" json:"uptime,omitempty"`
	UptimeEstimate int          `bson:"uptimeEstimate,omitempty" json:"uptimeEstimate,omitempty"`
	UptimeMillis   int          `bson:"uptimeMillis,omitempty" json:"uptimeMillis,omitempty"`
	Version        string       `bson:"version,omitempty" json:"version,omitempty"`
	Asserts        *Asserts     `bson:"asserts,omitempty" json:"asserts,omitempty"`
	Connections    *Connections `bson:"connections,omitempty" json:"connections,omitempty"`
	GlobalLock     *GlobalLock  `bson:"globalLock,omitempty" json:"globalLock,omitempty"`
	Network        *Network     `bson:"network,omitempty" json:"network,omitempty"`
	OPCounters     *OPCounters  `bson:"opcounters,omitempty" json:"opcounters,omitempty"`
	OPCountersRepl *OPCounters  `bson:"opcountersRepl,omitempty" json:"opcountersRepl,omitempty"`
	OK             int          `bson:"ok" json:"ok"`
	Error          string       `json:"error,omitempty"`
	// Repl           ServerRepl   `bson:"repl,omitempty" json:"repl,omitempty"`
}

// Asserts ...
type Asserts struct {
	Msg       int `bson:"msg" json:"msg"`
	Regular   int `bson:"regular" json:"regular"`
	Rollovers int `bson:"rollovers" json:"rollovers"`
	User      int `bson:"user" json:"user"`
	Warning   int `bson:"warning" json:"warning"`
}

// Connections ...
type Connections struct {
	Available    int `bson:"available" json:"available"`
	Current      int `bson:"current" json:"current"`
	TotalCreated int `bson:"totalCreated" json:"totalCreated"`
}

// GlobalLock ...
type GlobalLock struct {
	ActiveClients LockType `bson:"activeClients" json:"activeClients"`
	CurrentQueue  LockType `bson:"currentQueue" json:"currentQueue"`
	TotalTime     int      `bson:"totalTime" json:"totalTime"`
}

// LockType ...
type LockType struct {
	Readers int `bson:"readers" json:"readers"`
	Total   int `bson:"total" json:"total"`
	Writers int `bson:"writers" json:"writers"`
}

// Network ...
type Network struct {
	BytesIn     int `bson:"bytesIn" json:"bytesIn"`
	BytesOut    int `bson:"bytesOut" json:"bytesOut"`
	NumRequests int `bson:"numRequests" json:"numRequests"`
}

// OPCounters ...
type OPCounters struct {
	Command int `bson:"command" json:"command"`
	Delete  int `bson:"delete" json:"delete"`
	GetMore int `bson:"getmore" json:"getmore"`
	Insert  int `bson:"insert" json:"insert"`
	Query   int `bson:"query" json:"query"`
	Update  int `bson:"update" json:"update"`
}

// ServerRepl ...
type ServerRepl struct {
	SetName    string            `bson:"setName" json:"setName"`
	SetVersion int               `bson:"setVersion" json:"setVersion"`
	Me         string            `bson:"me" json:"me"`
	Primary    string            `bson:"primary" json:"primary"`
	IsMaster   bool              `bson:"isMaster" json:"isMaster"`
	Secondary  bool              `bson:"secondary" json:"secondary"`
	ElectionID bson.ObjectId     `bson:"electionId" json:"electionId"`
	RbID       int               `bson:"rbid" json:"rbid"`
	Arbiters   []string          `bson:"arbiters,omitempty" json:"arbiters,omitempty"`
	Hosts      []string          `bson:"hosts" json:"hosts"`
	Passives   []string          `bson:"passives,omitempty" json:"passives,omitempty"`
	Slaves     []ServerReplSlave `bson:"slaves" json:"slaves"`
	Tags       map[string]string `bson:"tags,omitempty" json:"tags,omitempty"`
}

// ServerReplSlave ...
type ServerReplSlave struct {
	Host     string        `bson:"host" json:"host"`
	MemberID int           `bson:"memberId" json:"memberId"`
	OPTime   int           `bson:"optime" json:"optime"`
	RID      bson.ObjectId `bson:"rid" json:"rid"`
}
