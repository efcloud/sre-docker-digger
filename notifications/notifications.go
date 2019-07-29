package notifications

import (
	"github.com/urfave/cli"
)

// Notifier interface definition
type Notifier interface {
	FireEvent(c *cli.Context, notification Notification) (err error)
}

// Notification defintion of Notification object
type Notification struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}
