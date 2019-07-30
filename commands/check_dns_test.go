package commands

import (
	"digger/notifications"
	"errors"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/urfave/cli"
)

// Setting up mock objects for Datadog API
type datadogNotifierMock struct{}

// NewNotification constructor for Notification
func NewdDatadogNotifierMock() notifications.Notifier {
	return &datadogNotifierMock{}
}

var callToDatadogFireEvent int

// FireEvent ...
func (event *datadogNotifierMock) FireEvent(c *cli.Context, notification notifications.Notification) (err error) {

	fmt.Print("run: datadogNotifierMock.FireEvent() \n")

	callToDatadogFireEvent++

	return nil
}

// Setting up mock objects for DNS API
type MyDNSClientMock struct{}

//NewMyDNSClientMock constructor for Client interface
func NewMyDNSClientMock() DNSClient {
	return &MyDNSClientMock{}
}

func (client *MyDNSClientMock) Exchange(target string, server string) (msg *dns.Msg, rtt time.Duration, err error) {

	fmt.Print("execute MyDNSClientMock.Exchange() \n")

	m := dns.Msg{}
	t := time.Duration(5)

	return &m, t, nil
}

// Setting up mock objects for DNS API, that returns an error
type MyDNSClientKOMock struct{}

//NewMyDNSClientMock constructor for Client interface
func NewMyDNSClientKOMock() DNSClient {
	return &MyDNSClientKOMock{}
}

func (client *MyDNSClientKOMock) Exchange(target string, server string) (msg *dns.Msg, rtt time.Duration, err error) {

	fmt.Print("execute MyDNSClientMock.Exchange() \n")

	m := dns.Msg{}
	t := time.Duration(5)

	return &m, t, errors.New("Expected error")
}

// TestRealActionCheckNoDatadog function to test realActionCheck()
func TestRealActionCheckNoDatadog(t *testing.T) {

	set := flag.NewFlagSet("test", 0)
	set.String("foo", "foo", "doc")
	context := cli.NewContext(nil, set, nil)

	err := realActionCheck(context, NewMyDNSClientMock(), NewdDatadogNotifierMock())

	if err != nil {
		t.Errorf("realActionCheck() should not have returned an error: %v", err)
	}

	if callToDatadogFireEvent > 0 {
		t.Error("FireEvent() has been called when it should not")
	}
}

// TestRealActionCheckDatadog function to test realActionCheck when Datadog is configured and should not be called
func TestRealActionCheckDatadogNotCalled(t *testing.T) {

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("datadog-enable", "true", "doc")
	globalContext := cli.NewContext(nil, globalSet, nil)

	set := flag.NewFlagSet("test", 0)
	set.String("interval", "2s", "doc")
	set.String("count", "1", "doc")
	context := cli.NewContext(nil, set, globalContext)

	fmt.Printf("local context: %s \n global context %s\n", context.String("interval"), context.GlobalString("datadog-enable"))

	err := realActionCheck(context, NewMyDNSClientMock(), NewdDatadogNotifierMock())

	if err != nil {
		t.Errorf("realActionCheck() should not have returned an error: %v", err)
	}

	if callToDatadogFireEvent > 0 {
		t.Error("FireEvent() should not have been called")
	}
}

// TestRealActionCheckDatadog function to test realActionCheck when Datadog is configured and should be called
func TestRealActionCheckDatadogCalled(t *testing.T) {

	globalSet := flag.NewFlagSet("test", 0)
	globalSet.String("datadog-enable", "true", "doc")
	globalContext := cli.NewContext(nil, globalSet, nil)

	set := flag.NewFlagSet("test", 0)
	set.String("interval", "2s", "doc")
	set.String("count", "1", "doc")
	context := cli.NewContext(nil, set, globalContext)

	fmt.Printf("local context: %s \n global context %s\n", context.String("interval"), context.GlobalString("datadog-enable"))

	err := realActionCheck(context, NewMyDNSClientKOMock(), NewdDatadogNotifierMock())

	if err != nil {
		t.Errorf("realActionCheck() should not have returned an error: %v", err)
	}

	if callToDatadogFireEvent < 1 {
		t.Error("FireEvent() should have been called")
	}
}