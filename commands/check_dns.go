package commands

import (
	"digger/dd"
	"digger/notifications"
	"strconv"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// DNSClient interface
type DNSClient interface {
	Exchange(target string, server string) (msg *dns.Msg, rtt time.Duration, err error)
}

// MyDNSClient struct
type MyDNSClient struct{}

//NewMyDNSClient constructor for Client interface
func NewMyDNSClient() DNSClient {
	return &MyDNSClient{}
}

// Exchange real implementation of Exchange()
func (client *MyDNSClient) Exchange(target string, server string) (msg *dns.Msg, rtt time.Duration, err error) {

	c := dns.Client{}
	m := dns.Msg{}

	m.SetQuestion(target+".", dns.TypeA)
	r, t, err := c.Exchange(&m, server+":53")
	return r, t, err
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

// CheckDNSCmd command definition
var CheckDNSCmd = cli.Command{
	Name:   "check-dns",
	Usage:  "Check connectivity.",
	Action: actionCheck,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "dns-server",
			Usage:  "DNS Server",
			EnvVar: "DNS_SERVER",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "FQDN to look for",
			EnvVar: "TARGET",
		},
		cli.StringFlag{
			Name:   "timeout",
			Usage:  "Timeout in seconds of the DNS request",
			EnvVar: "TIMEOUT",
			Value:  "5",
		},
		cli.StringFlag{
			Name:   "interval",
			Usage:  "Interval between 2 checks. Format is a number and a time suffix(eg: 1m, 5s, 10h).",
			EnvVar: "INTERVAL",
			Value:  "60s",
		},
		cli.StringFlag{
			Name:   "count",
			Usage:  "Number of check to run.",
			EnvVar: "COUNT",
			Value:  "0",
		},
	},
}

// actionCheck placeholder function
func actionCheck(c *cli.Context) (err error) {
	return realActionCheck(c, NewMyDNSClient(), dd.NewNotification())
}

func realActionCheck(c *cli.Context, client DNSClient, datadogNotifier notifications.Notifier) (err error) {
	if c.GlobalString("datadog-enable") == "true" {
		return runLoop(c, client, datadogNotifier)
	}

	return nil
}

// runLoop run the loop that check if the network is up
func runLoop(c *cli.Context, client DNSClient, notifier notifications.Notifier) (err error) {

	intervalDuration, err := time.ParseDuration(c.String("interval"))

	if err != nil {
		log.Error("Interval is not correct it should be a duration in the following format '60s' or '5m'")

		return err
	}

	count, _ := strconv.Atoi(c.String("count"))

	if count == 0 {
		for {
			err := runDNSCheck(c, client, notifier)
			if err != nil {
				log.Errorf("Erorr occured during DNS check: %v", err)
			}
			time.Sleep(intervalDuration)
		}
	} else {
		i := 0
		for i < count {
			err := runDNSCheck(c, client, notifier)
			if err != nil {
				log.Errorf("Erorr occured during DNS check: %v", err)
			}

			time.Sleep(intervalDuration)
			i++
		}
	}

	return nil
}

func runDNSCheck(c *cli.Context, client DNSClient, notifier notifications.Notifier) (err error) {

	_, t, err := client.Exchange(c.String("target"), c.String("dns-server"))

	log.Infof("Latency is %s", t)
	if err != nil {

		log.Error("Not able to reach remote DNS server, ", err)
		notification := notifications.Notification{
			Title: "Connectivity Issue",
			Text:  "Remote peer " + c.String("dns-server") + " is not reachable",
		}
		err = notifier.FireEvent(c, notification)

		if err != nil {
			log.Error("An error occured when sending an event to Datadog")
			return err
		}

		log.Info("Datadog event has been sent")

	}

	return nil
}
