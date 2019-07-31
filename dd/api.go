package dd

import (
	"bufio"
	"digger/notifications"
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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

	apiKey, appKey, err := getCredentials(c)
	if err != nil {
		return err
	}

	client := datadog.NewClient(apiKey, appKey)

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

func getCredentials(c *cli.Context) (apiKey string, appKey string, err error) {

	if (len(c.GlobalString("dd-api-key")) > 0) && (len(c.GlobalString("dd-app-key")) > 0) {
		return c.GlobalString("dd-api-key"), c.GlobalString("dd-app-key"), nil
	}

	if len(c.GlobalString("dd-creds-file")) > 0 {

		var (
			datadogAPIKey string
			datadogAppKey string
		)

		file, err := os.Open(c.GlobalString("dd-creds-file"))
		if err != nil {
			return "", "", err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := strings.Split(scanner.Text(), "=")
			switch s[0] {
			case "DATADOG_API_KEY":
				datadogAPIKey = s[1]
			case "DATADOG_APP_KEY":
				datadogAppKey = s[1]
			default:
				log.Warningf("Unknow key '%s' in Datadog credentials file", s[0])
			}
		}

		return datadogAPIKey, datadogAppKey, nil

	}

	return "", "", errors.New("Unable to locate Datadog credentials")
}
