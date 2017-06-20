package model

import (
	"time"
)

// ReplStatus ...
type ReplStatus struct {
	Date        *time.Time          `bson:"date,omitempty" json:"date,omitempty"`
	Members     []*ReplMember       `bson:"members,omitempty" json:"-"`
	JSONMembers map[int]*ReplMember `bson:"-" json:"members,omitempty"`
	MyState     *int                `bson:"myState,omitempty" json:"myState,omitempty"`
	Set         string              `bson:"set,omitempty" json:"set,omitempty"`
	SyncingTo   string              `bson:"syncingTo,omitempty" json:"syncingTo,omitempty"`
	OK          int                 `bson:"ok" json:"ok"`
	Error       string              `json:"error,omitempty"`
}

// ReplMember ...
type ReplMember struct {
	ID                   int        `bson:"_id" json:"_id"`
	ConfigVersion        int        `bson:"configVersion" json:"configVersion"`
	ElectionDate         *time.Time `bson:"electionDate,omitempty" json:"electionDate,omitempty"`
	ElectionTime         int        `bson:"electionTime,omitempty" json:"electionTime,omitempty"`
	Health               int        `bson:"health,omitempty" json:"health,omitempty"`
	LastHearbeat         *time.Time `bson:"lastHeartbeat,omitempty" json:"lastHeartbeat,omitempty"`
	LastHearbeatRecv     *time.Time `bson:"lastHeartbeatRecv,omitempty" json:"lastHeartbeatRecv,omitempty"`
	LastHeartbeatMessage string     `bson:"lastHeartbeatMessage,omitempty" json:"lastHeartbeatMessage,omitempty"`
	Name                 string     `bson:"name" json:"name"`
	OPTime               int64      `bson:"optime" json:"optime"`
	OPTimeDate           *time.Time `bson:"optimeDate" json:"optimeDate"`
	PingMS               int        `bson:"pingMs" json:"pingMs"`
	Self                 bool       `bson:"self,omitempty" json:"self,omitempty"`
	State                int        `bson:"state" json:"state"`
	StateStr             string     `bson:"stateStr" json:"stateStr"`
	Uptime               int        `bson:"uptime" json:"uptime"`
}
