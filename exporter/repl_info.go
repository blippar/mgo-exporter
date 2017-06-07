package exporter

import (
	"time"

	"github.com/blippar/mgo-exporter/mongo/model"
)

func getNodeReplInfo(repl *model.ReplStatus) *ReplicaSetInfo {

	if repl == nil || repl.Error != "" {
		return nil
	}

	members := repl.Members

	var me model.ReplMember
	for _, mb := range members {
		if mb.Self == true {
			me = mb
		}
	}

	var biggestLag int64
	oplogDiff := make(map[int]int64)
	for _, mb := range members {
		if mb.Self == true {
			continue
		}
		diff := mb.OPTimeDate.Sub(*me.OPTimeDate)
		lag := diff.Nanoseconds() / time.Second.Nanoseconds()
		if lag > biggestLag {
			biggestLag = lag
		}
		oplogDiff[mb.ID] = lag
	}

	return &ReplicaSetInfo{
		Set:           repl.Set,
		Name:          me.Name,
		State:         me.State,
		StateStr:      me.StateStr,
		OPTime:        me.OPTime,
		OPTimeDate:    me.OPTimeDate,
		Uptime:        me.Uptime,
		ConfigVersion: me.ConfigVersion,
		OPTimeDiff:    oplogDiff,
		OPTimeLag:     biggestLag,
	}
}
