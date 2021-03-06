package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	arg "github.com/alexflint/go-arg"
	log "github.com/sirupsen/logrus"

	"github.com/blippar/mgo-exporter/exporter"
	"github.com/blippar/mgo-exporter/forwarder"
	_ "github.com/blippar/mgo-exporter/forwarder/file"
	_ "github.com/blippar/mgo-exporter/forwarder/logstash"
)

// Version is the software version injected at build time
var Version = "unknown"

const (
	defaultForwarderURI = "file:///dev/stdout?pretty"
	defaultLogFile      = "/dev/stderr"
)

type cliArgs struct {
	MongoDB   string   `arg:"positional,help:Mongo URI for the node to connect to"`
	Database  []string `arg:"positional,separate,help:database name to monitor"`
	Repl      string   `arg:"-r,help:replicaSet name to monitor [env: MGOEXPORT_REPL],env:MGOEXPORT_REPL"`
	Forwarder string   `arg:"-f,help:forwarder URI to send messages to [env: MGOEXPORT_FORWARDER],env:MGOEXPORT_FORWARDER"`
	Logfile   string   `arg:"-l,help:file to output logs to [env: MGOEXPORT_LOGFILE],env:MGOEXPORT_LOGFILE"`
	Verbose   bool     `arg:"-v,help:enable a more verbose logging"`
	Quiet     bool     `arg:"-q,help:enable quieter logging"`
}

func (cliArgs) Version() string {
	return "mgo-exporter " + Version
}

func main() {

	args := &cliArgs{
		Forwarder: defaultForwarderURI,
		Logfile:   defaultLogFile,
	}
	arg.MustParse(args)

	// Init Logger
	if args.Logfile != defaultLogFile {
		f, err := os.OpenFile(args.Logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithError(err).
				WithField("logfile", args.Logfile).
				Fatal("initLoggerFileError")
		}
		defer func() {
			f.Sync()
			f.Close()
		}()
		log.SetOutput(f)
	}
	if args.Quiet && args.Verbose {
		log.WithField("description", "Multiple log level requested, will be set to default").
			Warn("initLoggerLevel")
	} else if args.Quiet {
		log.SetLevel(log.WarnLevel)
	} else if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// Init Forwarder
	fwd, err := forwarder.Factory.NewForwarder(args.Forwarder)
	if err != nil {
		log.WithError(err).WithField("uri", args.Forwarder).Fatal("initForwarderError")
	}
	log.WithField("uri", args.Forwarder).Info("initForwarderDone")

	// Init MongoDB Connection
	exporter, err := exporter.NewMongoStatsExporter(args.MongoDB, args.Database, args.Repl, fwd, 10*time.Second)
	if err != nil {
		log.WithError(err).Fatal("initMongoExporterError")
	}
	log.WithFields(exporter.LogFields).Info("initMongoExporterDone")
	exporter.Start()

	// Initialize signal channel
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait & Cleanup
	select {
	case <-exporter.Wait():
		log.WithField("description", "stopped due to a previous error").Fatal("stopped")
	case <-sigs:
		log.WithField("description", "stopped by signal").Info("stopped")
		exporter.Stop()
	}
	exporter.Close()
	fwd.Close()
}
