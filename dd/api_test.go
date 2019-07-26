package dd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/zorkian/go-datadog-api"
)

const (
	eventPayload = (`{
		"title": "My event",
		"text": "foobar"
	}
`)
)

// TestPostEventReturnCodeKO test function for PostEvent() handle correctly return code.
func TestPostEventReturnCodeOK(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	os.Setenv("DATADOG_HOST", ts.URL)

	ddEvent := datadog.Event{}

	_ = json.Unmarshal([]byte(eventPayload), &ddEvent)

	_, err := PostEvent("foo", "bar", &ddEvent)

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

	ddEvent := datadog.Event{}

	_ = json.Unmarshal([]byte(eventPayload), &ddEvent)

	_, err := PostEvent("foo", "bar", &ddEvent)

	if err == nil {
		t.Errorf("PostEvent() should have returned an error")
	}
}
