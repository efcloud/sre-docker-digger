package main

import (
	"testing"

	"github.com/urfave/cli"
)

// TestRunCli to ensure that the cli app is build as we expect
func TestRunCli(t *testing.T) {

	var (
		ddAPIKeyFlag      bool
		ddAppKeyFlag      bool
		ddCredsFileFlag   bool
		ddTagFlag         bool
		datadogEnableFlag bool
		checkCommand      bool
	)
	c := runCli()

	if c.Name != "digger" {
		t.Errorf("Expected application name to be 'digger', got '%v'", c.Name)
	}

	if c.HelpName != "CLI tool to check network connectivity." {
		t.Errorf("Expected application help to be 'CLI tool to check network connectivity.', got '%v'", c.Name)
	}

	for _, command := range c.Commands {
		switch command.Name {
		case "check-dns":
			checkCommand = true
		default:
			t.Errorf("Unexpected command '%s'", command.Name)
		}
	}

	if !checkCommand {
		t.Errorf("Command 'check-dns' is not defined")
	}

	for _, flag := range c.Flags {
		if f, ok := flag.(cli.StringFlag); ok {

			switch f.Name {
			case "dd-app-key":
				ddAppKeyFlag = true
			case "dd-api-key":
				ddAPIKeyFlag = true
			case "dd-creds-file":
				ddCredsFileFlag = true
			case "dd-tags":
				ddTagFlag = true
			case "datadog-enable":
				datadogEnableFlag = true

			default:
				t.Errorf("Unexpected flag '%s'", f.Name)
			}

		}
	}

	if !ddAPIKeyFlag {
		t.Errorf("Flag 'dd-api-key' is not defined")
	}

	if !ddAppKeyFlag {
		t.Errorf("Flag 'dd-app-key' is not defined")
	}

	if !ddCredsFileFlag {
		t.Errorf("Flag 'dd-creds-file' is not defined")
	}

	if !ddTagFlag {
		t.Errorf("Flag 'dd-tags' is not defined")
	}

	if !datadogEnableFlag {
		t.Errorf("Flag 'datadog-enable' is not defined")
	}

}
