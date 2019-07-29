package dd

import (
	"digger/notifications"

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

	ddEvent := datadog.Event{
		Title: &notification.Title,
		Text:  &notification.Text,
		// TODO: Need to use a flag to pass tags
		Tags:      []string{"environment:dev", "team:sre"},
		AlertType: &alertType,
		Priority:  &priority,
	}

	_, err = client.PostEvent(&ddEvent)

	if err != nil {
		return err
	}

	return nil
}
