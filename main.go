package main

import (
	"digger/commands"
	"os"

	"github.com/urfave/cli"
)

var version string

// main function
func main() {
	runCli().Run(os.Args)
}

// runCli build the CLI tool
func runCli() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "digger"
	app.HelpName = "CLI tool to check network connectivity."
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "datadog-enable",
			Usage:  "Enable or not Datadog notification",
			EnvVar: "DATADOG_ENABLE",
			Value:  "true",
		},
		cli.StringFlag{
			Name:   "dd-api-key",
			Usage:  "Datadog API key",
			EnvVar: "DATADOG_API_KEY",
		},
		cli.StringFlag{
			Name:   "dd-app-key",
			Usage:  "Datadog Application key",
			EnvVar: "DATADOG_APP_KEY",
		},
		cli.StringFlag{
			Name:   "dd-tags",
			Usage:  "Datadog tags, tags must be seperated by ','. For instance 'mytag1, key:value'.",
			EnvVar: "DATADOG_TAGS",
		},
	}
	app.Commands = []cli.Command{
		commands.CheckDNSCmd,
	}

	return
}
