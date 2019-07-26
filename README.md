# Digger

[![Build Status](https://drone.eu-west-1.edtech.sre.ef-cloud.io/api/badges/efcloud/sre-docker-digger/status.svg)](https://drone.eu-west-1.edtech.sre.ef-cloud.io/efcloud/sre-docker-digger) [![Go Report Card](https://goreportcard.com/badge/github.com/refcloud/sre-docker-digger)](https://goreportcard.com/report/efcloud/sre-docker-digger/dogsitter)
[![codecov](https://codecov.io/gh/efcloud/sre-docker-digger/branch/master/graph/badge.svg)](https://codecov.io/gh/efcloud/sre-docker-digger)



## Description
This is a repository for a small tool that will check if a connectivity is up by doing a DNS query.
In case of failure it will trigger an event to Datadog.

## Environment Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| DATADOG_HOST | URL of Datadog API endpoint | string | `https://api.datadoghq.eu` | no |
