# Digger

## Description
This is a repository for a small tool that will check if a connectivity is up by doing a DNS query.
In case of failure it will trigger an event to Datadog.

## Environment Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| DATADOG_HOST | URL of Datadog API endpoint | string | `https://api.datadoghq.eu` | no |
