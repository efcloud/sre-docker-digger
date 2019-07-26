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
			Name:   "enable-datadog",
			Usage:  "Enable or not Datadog notification",
			EnvVar: "DATADOG_API_KEY",
			Value:  "true",
		},
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "Datadog API key",
			EnvVar: "DATADOG_API_KEY",
		},
		cli.StringFlag{
			Name:   "app-key",
			Usage:  "Datadog Application key",
			EnvVar: "DATADOG_APP_KEY",
		},
	}
	app.Commands = []cli.Command{
		commands.CheckCmd,
	}

	return
}
