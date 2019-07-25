package dd

import (
	log "github.com/sirupsen/logrus"
	"github.com/zorkian/go-datadog-api"
)

// PostEvent fire an event to Datadog
func PostEvent(apiKey string, appKey string, event datadog.Event) (err error) {
	client := datadog.NewClient(apiKey, appKey)
	log.Info("dd client ", client)
	return nil
}
