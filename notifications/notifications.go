package notifications

import (
	"digger/dd"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/zorkian/go-datadog-api"
)

// Notification interface definition
type Notification interface {
	FireEvent(c *cli.Context) (err error)
	fireDatadogEvent(apiKey string, appKey string) (err error)
}

// DiggerNotification defintion of Notification object
type DiggerNotification struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// NewDiggerNotification constructor for DiggerNotification
func NewDiggerNotification(title string, text string) Notification {
	return &DiggerNotification{Title: title, Text: text}
}

// FireEvent function to firevent using the enabled Notificaton service
func (event *DiggerNotification) FireEvent(c *cli.Context) (err error) {

	if c.GlobalString("datadog-enable") == "true" {
		err = event.fireDatadogEvent(c.GlobalString("dd-api-key"), c.GlobalString("dd-app-key"))
		if err != nil {
			return err
		}
	}
	return nil
}

// fireDatadogEvent sent a event to Datadog
func (event *DiggerNotification) fireDatadogEvent(apiKey string, c string) (err error) {

	alertType := "error"
	priority := "normal"

	ddEvent := datadog.Event{
		Title:     &event.Title,
		Text:      &event.Text,
		Tags:      []string{"environment:dev", "team:sre"},
		AlertType: &alertType,
		Priority:  &priority,
	}

	postedEvent, err := dd.PostEvent(apiKey, apiKey, &ddEvent)

	if err != nil {
		return err
	}

	log.Infof("Event posted id %d", *postedEvent.Id)
	return nil

}
