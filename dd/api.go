package dd

import (
	"github.com/zorkian/go-datadog-api"
)

// PostEvent fire an event to Datadog
func PostEvent(apiKey string, appKey string, event *datadog.Event) (postedEvent *datadog.Event, err error) {
	client := datadog.NewClient(apiKey, appKey)

	postedEvent, err = client.PostEvent(event)

	if err != nil {
		return nil, err
	}

	return postedEvent, nil
}
