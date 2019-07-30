package dd

import (
	"digger/notifications"
	"flag"
	"io/ioutil"
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

// TestFireEventReturnCodeOK test function for FireEvent() returns an error when creds are absent.
func TestFireEventKO(t *testing.T) {

	event := NewNotification()

	notification := notifications.Notification{
		Title: "Test",
		Text:  "Test text",
	}

	set := flag.NewFlagSet("test", 0)
	context := cli.NewContext(nil, set, nil)

	err := event.FireEvent(context, notification)

	if err == nil {
		t.Errorf("FireEvent() should have returned an error")
	}
}

// TestFireEventReturnCodeOK test function for FireEvent() handle correctly return code.
func TestFireEventReturnCodeOK(t *testing.T) {

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
		t.Errorf("FireEvent() should not have returned an error: %v", err)
	}
}

// TestFireEventReturnCodeKO test function for FireEvent() handle correctly return code.
func TestFireEventReturnCodeKO(t *testing.T) {

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
		t.Errorf("FireEvent() should have returned an error")
	}
}

// TestParseTag function to test parseTag()
func TestParseTag(t *testing.T) {

	expectedTags := []string{"foo:bar", "titi:toto"}

	returnedTags := parseTag("foo:bar,titi:toto")

	if len(returnedTags) != len(expectedTags) {
		t.Errorf("parseTag() should have returned a list of %d elements, got %d", len(expectedTags), len(returnedTags))
	}

	for i, v := range expectedTags {
		if v != returnedTags[i] {
			t.Errorf("parseTag() should have returned a list  with the following element %s, got %s", expectedTags[i], returnedTags[i])
		}
	}
}

// TestGetCredentialsNoCreds function to getCredentials() when no creds are passed
func TestGetCredentialsNoCreds(t *testing.T) {

	globalSet := flag.NewFlagSet("test", 0)
	globalContext := cli.NewContext(nil, globalSet, nil)

	context := cli.NewContext(nil, nil, globalContext)

	_, _, err := getCredentials(context)

	if err == nil {
		t.Errorf("getCredentials() should have retuned and error")
	}

}

// TestGetCredentials function to getCredentials()
func TestGetCredentials(t *testing.T) {

	expectedAPIKey := "aaa"
	expectedAppKey := "bbb"

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("dd-api-key", expectedAPIKey, "doc")
	globalSet.String("dd-app-key", expectedAppKey, "doc")
	globalContext := cli.NewContext(nil, globalSet, nil)

	context := cli.NewContext(nil, nil, globalContext)

	apiKey, appKey, err := getCredentials(context)

	if err != nil {
		t.Errorf("getCredentials() should not have retuned and error, %v", err)
	}

	if apiKey != expectedAPIKey {
		t.Errorf("getCredentials() should have this API KEY %s, got %s", expectedAPIKey, apiKey)
	}

	if appKey != expectedAppKey {
		t.Errorf("getCredentials() should have this APP KEY %s, got %s", expectedAppKey, appKey)
	}

}

// TestGetCredentialsFile function to getCredentials() when a credentials file is used
func TestGetCredentialsFile(t *testing.T) {

	expectedAPIKey := "aaa"
	expectedAppKey := "bbb"
	credsContent := (`DATADOG_APP_KEY=bbb
DATADOG_API_KEY=aaa`)
	credsFile := "/tmp/dd-creds.file"

	err := ioutil.WriteFile(credsFile, []byte(credsContent), 0644)

	if err != nil {
		t.Errorf("Error when building %s: %v", credsFile, err)
	}

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("dd-creds-file", credsFile, "doc")
	globalContext := cli.NewContext(nil, globalSet, nil)

	context := cli.NewContext(nil, nil, globalContext)

	apiKey, appKey, err := getCredentials(context)

	if err != nil {
		t.Errorf("getCredentials() should not have retuned and error, %v", err)
	}

	if apiKey != expectedAPIKey {
		t.Errorf("getCredentials() should have this API KEY %s, got %s", expectedAPIKey, apiKey)
	}

	if appKey != expectedAppKey {
		t.Errorf("getCredentials() should have this APP KEY %s, got %s", expectedAppKey, appKey)
	}

	// Cleaning up
	_ = os.Remove(credsFile)

}

// TestGetCredentialsFile function to getCredentials() when a credentials file is used
func TestGetCredentialsFileNotExist(t *testing.T) {

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("dd-creds-file", "none", "doc")
	globalContext := cli.NewContext(nil, globalSet, nil)

	context := cli.NewContext(nil, nil, globalContext)

	_, _, err := getCredentials(context)

	if err == nil {
		t.Errorf("getCredentials() should have retuned and error")
	}

}
