package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"

	"github.com/blippar/mgo-exporter/exporter"
	"github.com/blippar/mgo-exporter/logstash"
)

const (
	defaultLogstashHost = "127.0.0.1:2000"
)

func main() {

	var args struct {
		Host     string   `arg:"positional"`
		Port     int      `arg:"positional"`
		Database []string `arg:"-d,separate"`
		Logstash string   `arg:"-l"`
		Repl     string   `arg:"-r,help:Monitor Replica Set"`
		Verbose  bool     `arg:"-v"`
	}
	args.Logstash = defaultLogstashHost
	arg.MustParse(&args)

	log.SetHandler(text.New(os.Stderr))
	if args.Verbose {
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
	exporter, err := exporter.NewMongoStatsExporter(args.Host, args.Port, args.Database, args.Repl, forwarder, 10*time.Second)
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
