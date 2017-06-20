package exporter

import (
	"time"

	log "github.com/apex/log"
	mgo "gopkg.in/mgo.v2"

	"github.com/blippar/mgo-exporter/forwarder"
	"github.com/blippar/mgo-exporter/mongo"
)

// MongoStatsExporter ...
type MongoStatsExporter struct {
	session   *mgo.Session
	info      ServerInfo
	databases []string
	forwarder forwarder.Forwarder
	every     time.Duration
	doneCh    chan struct{}
	stopCh    chan struct{}
	logFields log.Fields
}

// NewMongoStatsExporter ...
func NewMongoStatsExporter(connURI string, databases []string, repl string, fwd forwarder.Forwarder, every time.Duration) (*MongoStatsExporter, error) {

	// Set MongoDB driver configuration
	dialInfo, err := mgo.ParseURL(connURI)
	if err != nil {
		return nil, err
	}
	dialInfo.Timeout = 5 * time.Second
	dialInfo.PoolLimit = 512
	dialInfo.FailFast = true
	dialInfo.Direct = true
	dialInfo.ReplicaSetName = repl
	// dialInfo.Username = ""
	// dialInfo.Password = ""

	// Connect to MongoDB
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)

	// Default LogFields
	logFields := log.Fields{
		"host": dialInfo.Addrs[0],
	}
	if repl != "" {
		logFields["repl"] = repl
	}

	// Return new object
	return &MongoStatsExporter{
		session: session,
		info: ServerInfo{
			Host:       dialInfo.Addrs[0],
			ReplicaSet: repl,
		},
		databases: databases,
		forwarder: fwd,
		every:     every,
		logFields: logFields,
	}, nil
}

func (e *MongoStatsExporter) isConnected() error {
	// Ping to check if the connection is still alive
	return e.session.Ping()
}

func (e *MongoStatsExporter) reconnect() error {
	// Refresh to try to reconnect
	e.session.Refresh()
	// Check if connection is now alive
	return e.isConnected()
}

func (e *MongoStatsExporter) export(exportTime time.Time) []*Message {

	var err error

	sMsg := make([]*Message, 0, 1+len(e.databases))
	mongoInfo := e.info
	mongoInfo.Connected = true

	m := &Message{
		Time:  exportTime,
		Mongo: &mongoInfo,
		Type:  "serverStatus",
	}
	sMsg = append(sMsg, m)

	log.WithFields(e.logFields).Info("mongoExport")
	// Export Server status then if sucessful try to excract other info
	m.ServerStatus, err = mongo.GetServerStatus(e.session)
	if err != nil {
		log.WithFields(e.logFields).WithError(err).Error("getServerStatusError")
		return sMsg
	}

	// Export ReplicaSet status
	if e.info.ReplicaSet != "" {
		m.Repl, err = mongo.GetReplStatus(e.session)
		if err != nil {
			log.WithFields(e.logFields).WithError(err).Warn("getReplStatusError")
		}
		m.NodeReplInfo = getNodeReplInfo(m.Repl)
	}

	// Export DB stats
	for _, db := range e.databases {

		log.WithFields(e.logFields).WithField("db", db).Info("mongoDBExport")

		dbStats, err := mongo.GetDBStats(e.session, db)
		if err != nil {
			log.WithFields(e.logFields).WithError(err).Warn("getDBStatsError")
		}
		dmsg := &Message{
			Time:    exportTime,
			Mongo:   &mongoInfo,
			DBStats: dbStats,
			Type:    "dbStats",
		}
		sMsg = append(sMsg, dmsg)
	}

	return sMsg
}

func (e *MongoStatsExporter) run() {

	ticker := time.NewTicker(e.every)
	defer ticker.Stop()
	defer close(e.doneCh)

	for {
		select {
		case <-ticker.C:

			var errCo error
			var messages []*Message
			ctime := time.Now()

			// Try to check if connected and try to reconnect before exporting
			if err := e.isConnected(); err != nil {
				errCo = e.reconnect()
			}

			// If we can't connect, create an error message
			if errCo != nil {

				log.WithFields(e.logFields).WithError(errCo).Warn("mongoConError")

				mongoInfo := e.info
				mongoInfo.Connected = false
				mongoInfo.Error = errCo.Error()

				messages = []*Message{
					{
						Time:  ctime,
						Mongo: &mongoInfo,
						Type:  "serverStatus",
					},
				}

			} else { // Else export MongoDB data

				messages = e.export(ctime)

			}

			// Send Messages to Logstash
			for _, m := range messages {
				if err := e.forwarder.Send(m); err != nil {
					log.WithError(err).Error("logstashSendMessageError")
					return
				}
			}

		case <-e.stopCh:
			return
		}
	}
}

// Start ...
func (e *MongoStatsExporter) Start() chan struct{} {

	e.doneCh = make(chan struct{})
	e.stopCh = make(chan struct{})

	go e.run()
	return e.doneCh
}

// Stop ...
func (e *MongoStatsExporter) Stop() {

	select {
	case e.stopCh <- struct{}{}:
		<-e.doneCh
	case <-e.doneCh:
	}
}

// Wait ...
func (e *MongoStatsExporter) Wait() <-chan struct{} {
	return e.doneCh
}

// Close ...
func (e *MongoStatsExporter) Close() {
	e.session.Close()
}

// NOTE: export only if PRIMARY or SOLO
// NOTE: see about nested field
