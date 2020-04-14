package main

import (
	"github.com/Percona-Lab/mnogo_exporter/exporter"
	"github.com/alecthomas/kong"
)

// GlobalFlags has command line flags to configure the exporter
type GlobalFlags struct {
	//DSN     string `required:"true" help:"MongoDB connection URI" placeholder:"mongodb://user:pass@127.0.0.1:27017/admin?ssl=true"`
	DSN     string `help:"MongoDB connection URI" placeholder:"mongodb://user:pass@127.0.0.1:27017/admin?ssl=true"`
	Debug   bool   `short:"D" help:"Enable debug mode"`
	Version bool   `help:"Show version and exit"`
}

func main() {
	var opts GlobalFlags
	ctx := kong.Parse(&opts,
		kong.Name("mnogo_exporter"),
		kong.Description("MongoDB Prometheus exporter"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.0.1",
		})

	exporterOpts := &exporter.Opts{
		DSN: opts.DSN,
	}

	e, err := exporter.New(exporterOpts)
	if err != nil {
		panic(err)
	}
	_ = e
}
