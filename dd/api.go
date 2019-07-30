package dd

import (
	"digger/notifications"
	"strings"

	"github.com/urfave/cli"
	"github.com/zorkian/go-datadog-api"
)

type datadogNotifier struct{}

// NewNotification constructor for Notification
func NewNotification() notifications.Notifier {
	return &datadogNotifier{}
}

// FireEvent function to post an event in Datadog
func (event *datadogNotifier) FireEvent(c *cli.Context, notification notifications.Notification) (err error) {

	client := datadog.NewClient(c.GlobalString("dd-api-key"), c.GlobalString("dd-app-key"))

	alertType := "error"
	priority := "normal"
	tags := parseTag(c.GlobalString("dd-tags"))

	ddEvent := datadog.Event{
		Title:     &notification.Title,
		Text:      &notification.Text,
		Tags:      tags,
		AlertType: &alertType,
		Priority:  &priority,
	}

	_, err = client.PostEvent(&ddEvent)

	if err != nil {
		return err
	}

	return nil
}

func parseTag(s string) []string {

	tags := strings.Split(s, ",")

	return tags
}
