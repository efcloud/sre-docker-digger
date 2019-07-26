package notifications

import (
	"github.com/urfave/cli"
)

// Event ...
type Event struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// FireEvent ...
func FireEvent(c *cli.Context, event Event) (err error) {

	return nil
}
