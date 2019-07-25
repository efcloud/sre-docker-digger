package commands

import (
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

// CheckCmd ...
var CheckCmd = cli.Command{
	Name:  "check",
	Usage: "Check connectivity.",
	// Action: NewDelete().delete,
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
	},
}

// actionCheck placeholder function
func actionCheck(c *cli.Context) (err error) {
	return runLoop(c, NewMyDNSClient())
}

// runLoop ...
func runLoop(c *cli.Context, client DNSClient) (err error) {

	intervalDuration, err := time.ParseDuration(c.String("interval"))

	if err != nil {
		log.Error("Interval is correct it should be a duration in")
	}

	for {
		t, err := runTest(client, c.String("target"), c.String("dns-server"))
		log.Info("Latency is ", t)
		if err != nil {
			log.Error("Not able to reach remote DNS server, ", err)
			return err
		}
		time.Sleep(intervalDuration)
	}

}

func runTest(client DNSClient, target string, server string) (duration time.Duration, err error) {
	_, t, err := client.Exchange(target, server)

	return t, err
}
