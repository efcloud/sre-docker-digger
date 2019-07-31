# Digger

[![Build Status](https://drone.eu-west-1.edtech.sre.ef-cloud.io/api/badges/efcloud/sre-docker-digger/status.svg)](https://drone.eu-west-1.edtech.sre.ef-cloud.io/efcloud/sre-docker-digger) [![Go Report Card](https://goreportcard.com/badge/github.com/efcloud/sre-docker-digger)](https://goreportcard.com/report/efcloud/sre-docker-digger)
[![codecov](https://codecov.io/gh/efcloud/sre-docker-digger/branch/master/graph/badge.svg)](https://codecov.io/gh/efcloud/sre-docker-digger)



## Description
This is a repository for a small tool that will check if a connectivity is up by doing a DNS query.  
In case of failure it can be configured to fire an a notification.  
For now only Datadog is supported for notifications.

## Environment Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| DATADOG_HOST | URL of Datadog API endpoint | string | `https://api.datadoghq.eu` | no |

## Usage
```
$ digger
NAME:
   digger - A new cli application

USAGE:
   CLI tool to check network connectivity. [global options] command [command options] [arguments...]

COMMANDS:
     check-dns  Check connectivity.
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --datadog-enable value  Enable or not Datadog notification (default: "true") [$DATADOG_ENABLE]
   --dd-api-key value      Datadog API key [$DATADOG_API_KEY]
   --dd-app-key value      Datadog Application key [$DATADOG_APP_KEY]
   --dd-creds-file value   File containning Datadog credentials. [$DATADOG_CREDENTIALS_FILE]
   --dd-tags value         Datadog tags, tags must be seperated by ','. For instance 'mytag1, key:value'. [$DATADOG_TAGS]
   --help, -h              show help
   --version, -v           print the version
```

### Check-dns command
```
$ digger check-dns -h
NAME:
   CLI tool to check network connectivity. check-dns - Check connectivity.

USAGE:
   CLI tool to check network connectivity. check-dns [command options] [arguments...]

OPTIONS:
   --dns-server value  DNS Server [$DNS_SERVER]
   --target value      FQDN to look for [$TARGET]
   --timeout value     Timeout in seconds of the DNS request (default: "5") [$TIMEOUT]
   --interval value    Interval between 2 checks. Format is a number and a time suffix(eg: 1m, 5s, 10h). (default: "60s") [$INTERVAL]
   --count value       Number of check to run. '0' run the test forever. (default: "0") [$COUNT]
```

## Notifications

### Datadog
If you want digger to notify Datadog when the connectivity is not working set `--datadog-enable` to `true`.  

To pass Datadog credentials you can either use the CLI flags `--dd-api-key` and `--dd-app-key` or pass a file containing those keys using this flag `--dd-creds-file`.

The file needs to be like this (order does not matter):
```
DATADOG_API_KEY=XXX
DATADOG_APP_KEY=XXX
```
