package main

import (
	"log"
	"os"

	"github.com/Percona-Lab/mnogo_exporter/exporter"
	"gopkg.in/alecthomas/kingpin.v2"
)

// GlobalFlags has command line flags to configure the exporter
type GlobalFlags struct {
	DSN     string
	Debug   bool
	Version bool
}

func main() {
	var opts GlobalFlags

	app := kingpin.New("mnogo_exporter", "MongoDB metrics exporter")
	app.Flag("debug", "Enable debug mode.").BoolVar(&opts.Debug)
	app.Flag("mongodb.uri", "MongoDB connection string").StringVar(&opts.DSN)
	app.Flag("version", "Show version and exit").BoolVar(&opts.Version)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Cannot parse command line arguments: %s", err)
	}

	exporterOpts := &exporter.Opts{
		DSN: opts.DSN,
	}
	e, err := exporter.New(exporterOpts)
	if err != nil {
		panic(err)
	}
	_ = e
}
