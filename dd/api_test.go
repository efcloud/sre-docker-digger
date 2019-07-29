package dd

import (
	"digger/notifications"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/urfave/cli"
)

const (
	eventPayload = (`{
		"title": "My event",
		"text": "foobar"
	}
`)
)

// TestPostEventReturnCodeOK test function for PostEvent() handle correctly return code.
func TestPostEventReturnCodeOK(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	os.Setenv("DATADOG_HOST", ts.URL)

	event := NewNotification()

	notification := notifications.Notification{
		Title: "Test",
		Text:  "Test text",
	}

	set := flag.NewFlagSet("test", 0)
	set.String("dd-api-key", "foo", "doc")
	set.String("dd-app-key", "bar", "doc")
	context := cli.NewContext(nil, set, nil)

	err := event.FireEvent(context, notification)

	if err != nil {
		t.Errorf("PostEvent() should not have returned an error: %v", err)
	}
}

// TestPostEventReturnCodeKO test function for PostEvent() handle correctly return code.
func TestPostEventReturnCodeKO(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	os.Setenv("DATADOG_HOST", ts.URL)

	event := NewNotification()

	notification := notifications.Notification{
		Title: "Test",
		Text:  "Test text",
	}

	set := flag.NewFlagSet("test", 0)
	set.String("dd-api-key", "foo", "doc")
	set.String("dd-app-key", "bar", "doc")
	context := cli.NewContext(nil, set, nil)

	err := event.FireEvent(context, notification)

	if err == nil {
		t.Errorf("PostEvent() should have returned an error")
	}
}
