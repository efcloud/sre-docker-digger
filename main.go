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

	app.Commands = []cli.Command{
		commands.CheckCmd,
	}

	return
}
