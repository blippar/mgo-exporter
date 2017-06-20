package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/blippar/mgo-exporter/exporter"
	"github.com/blippar/mgo-exporter/logstash"
)

const (
	defaultLogstashHost = "127.0.0.1:2000"
)

func main() {

	var args struct {
		MongoDB  string   `arg:"positional,help:Mongo URI for the node to connect to"`
		Database []string `arg:"-d,separate,help:database name to monitor"`
		Repl     string   `arg:"-r,help:replicaSet name to monitor"`
		Logstash string   `arg:"-l,help:Logstash URI to send messages to"`
		Verbose  bool     `arg:"-v,help:enable a more verbose logging"`
		Quiet    bool     `arg:"-q,help:enable quieter logging"`
	}
	args.Logstash = defaultLogstashHost
	arg.MustParse(&args)

	log.SetHandler(text.New(os.Stderr))
	if args.Quiet {
		log.SetHandler(json.New(os.Stderr))
		log.SetLevel(log.WarnLevel)
	} else if args.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Init Logstash Forwarder
	forwarder, err := logstash.NewTCPForwarder(args.Logstash)
	if err != nil {
		log.WithError(err).Fatal("initLogstashForwarderError")
	}

	// Init MongoDB Connection
	exporter, err := exporter.NewMongoStatsExporter(args.MongoDB, args.Database, args.Repl, forwarder, 10*time.Second)
	if err != nil {
		log.WithError(err).Fatal("initMongoStatsExporterError")
	}
	exporter.Start()

	// Initialize signal channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait & Cleanup
	select {
	case <-exporter.Wait():
		log.WithField("description", "stopped due to a previous error").Info("stopped")
	case <-sigs:
		log.WithField("description", "stopped by signal").Info("stopped")
		exporter.Stop()
	}
	exporter.Close()
	forwarder.Close()
}
