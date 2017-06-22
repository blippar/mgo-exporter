package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ServerStatus ...
type ServerStatus struct {
	Host             string           `bson:"host,omitempty"             json:"host,omitempty"`
	Version          string           `bson:"version,omitempty"          json:"version,omitempty"`
	Process          string           `bson:"process,omitempty"          json:"process,omitempty"`
	PID              int              `bson:"pid,omitempty"              json:"pid,omitempty"`
	Uptime           int              `bson:"uptime,omitempty"           json:"uptime,omitempty"`
	UptimeEstimate   int              `bson:"uptimeEstimate,omitempty"   json:"uptimeEstimate,omitempty"`
	UptimeMillis     int              `bson:"uptimeMillis,omitempty"     json:"uptimeMillis,omitempty"`
	LocalTime        *time.Time       `bson:"localTime,omitempty"        json:"localTime,omitempty"`
	Asserts          *Asserts         `bson:"asserts,omitempty"          json:"asserts,omitempty"`
	Connections      *Connections     `bson:"connections,omitempty"      json:"connections,omitempty"`
	GlobalLock       *GlobalLock      `bson:"globalLock,omitempty"       json:"globalLock,omitempty"`
	Network          *Network         `bson:"network,omitempty"          json:"network,omitempty"`
	OPCounters       *OPCounters      `bson:"opcounters,omitempty"       json:"opcounters,omitempty"`
	OPCountersRepl   *OPCounters      `bson:"opcountersRepl,omitempty"   json:"opcountersRepl,omitempty"`
	Mem              *Memory          `bson:"mem,omitempty"              json:"mem,omitempty"`
	Dur              *Journaling      `bson:"dur,omitempty"              json:"-,omitempty"`
	ExtraInfo        *ExtraInfo       `bson:"extra_info,omitempty"       json:"-"`
	Locks            map[string]*Lock `bson:"locks,omitempty"            json:"-"`
	StorageEngine    *StorageEngine   `bson:"storageEngine,omitempty"    json:"-"`
	Metrics          *Metrics         `bson:"metrics,omitempty"          json:"-"`
	WiredTiger       *WiredTiger      `bson:"wiredTiger,omitempty"       json:"-"`
	Repl             *ServerRepl      `bson:"repl,omitempty"             json:"-"`
	WriteBacksQueued *bool            `bson:"writeBacksQueued,omitempty" json:"-"`
	OK               int              `bson:"ok"                         json:"ok"`
	Error            string           `bson:"-"                          json:"error,omitempty"`
}

// Asserts ...
type Asserts struct {
	Msg       int `bson:"msg"       json:"msg"`
	Regular   int `bson:"regular"   json:"regular"`
	Rollovers int `bson:"rollovers" json:"rollovers"`
	User      int `bson:"user"      json:"user"`
	Warning   int `bson:"warning"   json:"warning"`
}

// BackgroundFlushing ...
type BackgroundFlushing struct {
	Flushes      int        `bson:"flushes"       json:"flushes"`
	TotalMS      int        `bson:"total_ms"      json:"total_ms"`
	AverageMS    float64    `bson:"average_ms"    json:"average_ms"`
	LastMS       int        `bson:"last_ms"       json:"last_ms"`
	LastFinished *time.Time `bson:"last_finished" json:"last_finished"`
}

// Connections ...
type Connections struct {
	Available    int `bson:"available"    json:"available"`
	Current      int `bson:"current"      json:"current"`
	TotalCreated int `bson:"totalCreated" json:"totalCreated"`
}

// Journaling ...
type Journaling struct {
	Commits            int            `bson:"commits"            json:"commits"`
	JournaledMB        int            `bson:"journaledMB"        json:"journaledMB"`
	WriteToDataFilesMB int            `bson:"writeToDataFilesMB" json:"writeToDataFilesMB"`
	Compression        int            `bson:"compression"        json:"compression"`
	CommitsInWriteLock int            `bson:"commitsInWriteLock" json:"commitsInWriteLock"`
	EarlyCommits       int            `bson:"earlyCommits"       json:"earlyCommits"`
	TimeMs             JournalingTime `bson:"timeMs"             json:"timeMs"`
}

// JournalingTime ...
type JournalingTime struct {
	DT                 int `bson:"dt"                 json:"dt"`
	PrepLogBuffer      int `bson:"prepLogBuffer"      json:"prepLogBuffer"`
	WriteToJournal     int `bson:"writeToJournal"     json:"writeToJournal"`
	WriteToDataFiles   int `bson:"writeToDataFiles"   json:"writeToDataFiles"`
	RemapPrivateView   int `bson:"remapPrivateView"   json:"remapPrivateView"`
	Commits            int `bson:"commits"            json:"commits"`
	CommitsInWriteLock int `bson:"commitsInWriteLock" json:"commitsInWriteLock"`
}

// ExtraInfo ...
type ExtraInfo struct {
	Note           string `bson:"note"                       json:"note"`
	HeapUsageBytes int    `bson:"heap_usage_bytes,omitempty" json:"heap_usage_bytes,omitempty"`
	PageFaults     int    `bson:"page_faults"                json:"page_faults"`
}

// GlobalLock ...
type GlobalLock struct {
	ActiveClients GlobalLockType `bson:"activeClients" json:"activeClients"`
	CurrentQueue  GlobalLockType `bson:"currentQueue"  json:"currentQueue"`
	TotalTime     int            `bson:"totalTime"     json:"totalTime"`
}

// GlobalLockType ...
type GlobalLockType struct {
	Readers int `bson:"readers" json:"readers"`
	Total   int `bson:"total"   json:"total"`
	Writers int `bson:"writers" json:"writers"`
}

// Lock ...
type Lock struct {
	AcquireCount        *LockAcquire `bson:"acquireCount,omitempty"        json:"acquireCount,omitempty"`
	AcquireWaitCount    *LockAcquire `bson:"acquireWaitCount,omitempty"    json:"acquireWaitCount,omitempty"`
	TimeAcquiringMicros *LockAcquire `bson:"timeAcquiringMicros,omitempty" json:"timeAcquiringMicros,omitempty"`
	DeadlockCount       *LockAcquire `bson:"deadlockCount,omitempty"       json:"deadlockCount,omitempty"`
}

// LockAcquire ...
type LockAcquire struct {
	Shared          *int `bson:"R,omitempty" json:"R,omitempty"`
	Exclusive       *int `bson:"W,omitempty" json:"W,omitempty"`
	IntentShared    *int `bson:"r,omitempty" json:"r,omitempty"`
	IntentExclusive *int `bson:"w,omitempty" json:"w,omitempty"`
}

// Network ...
type Network struct {
	BytesIn     int `bson:"bytesIn"     json:"bytesIn"`
	BytesOut    int `bson:"bytesOut"    json:"bytesOut"`
	NumRequests int `bson:"numRequests" json:"numRequests"`
}

// OPCounters ...
type OPCounters struct {
	Command int `bson:"command" json:"command"`
	Delete  int `bson:"delete"  json:"delete"`
	GetMore int `bson:"getmore" json:"getmore"`
	Insert  int `bson:"insert"  json:"insert"`
	Query   int `bson:"query"   json:"query"`
	Update  int `bson:"update"  json:"update"`
}

// ServerRepl ...
type ServerRepl struct {
	SetName    string            `bson:"setName"            json:"setName"`
	SetVersion int               `bson:"setVersion"         json:"setVersion"`
	Me         string            `bson:"me"                 json:"me"`
	Primary    string            `bson:"primary"            json:"primary"`
	IsMaster   bool              `bson:"isMaster"           json:"isMaster"`
	Secondary  bool              `bson:"secondary"          json:"secondary"`
	ElectionID bson.ObjectId     `bson:"electionId"         json:"electionId"`
	RbID       int               `bson:"rbid"               json:"rbid"`
	Arbiters   []string          `bson:"arbiters,omitempty" json:"arbiters,omitempty"`
	Hosts      []string          `bson:"hosts"              json:"hosts"`
	Passives   []string          `bson:"passives,omitempty" json:"passives,omitempty"`
	Slaves     []ServerReplSlave `bson:"slaves"             json:"slaves"`
	Tags       map[string]string `bson:"tags,omitempty"     json:"tags,omitempty"`
}

// ServerReplSlave ...
type ServerReplSlave struct {
	Host     string        `bson:"host"     json:"host"`
	MemberID int           `bson:"memberId" json:"memberId"`
	OPTime   int           `bson:"optime"   json:"optime"`
	RID      bson.ObjectId `bson:"rid"      json:"rid"`
}

// Memory ...
type Memory struct {
	Bits              int `bson:"bits"              json:"bits"`
	Resident          int `bson:"resident"          json:"resident"`
	Virtual           int `bson:"virtual"           json:"virtual"`
	Supported         int `bson:"supported"         json:"supported"`
	Mapped            int `bson:"mapped"            json:"mapped"`
	MappedWithJournal int `bson:"mappedWithJournal" json:"mappedWithJournal"`
}

// StorageEngine ...
type StorageEngine struct {
	Name string `bson:"name" json:"name"`
}

// Metrics ...
type Metrics map[string]interface{}

// WiredTiger ...
type WiredTiger map[string]interface{}
