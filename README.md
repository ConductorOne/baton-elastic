# `baton-elastic` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-elastic.svg)](https://pkg.go.dev/github.com/conductorone/baton-elastic) ![main ci](https://github.com/conductorone/baton-elastic/actions/workflows/main.yaml/badge.svg)

`baton-elastic` is a connector for Elastic built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Elastic API to sync data about Elastic cloud organizations and users. Optionally it can also sync elasticsearch deployment roles and users.
Check out [Baton](https://github.com/conductorone/baton) to learn more about the project in general.

# Getting Started

## Prerequisites

- Access to the Elastic cloud.
- API key to access Elastic cloud API. You can create the key in Organization -> API keys
- By default the connector will sync only organizations and users from Elastic cloud. If you also want to sync users and roles from a specific deployment simply provide the `--deployment-endpoint` and `--deployment-api-key` flags. You can find your deployment endpoint in the top right corner of Integration page under 'Connection details' -> Elasticsearch endpoint. To create an API key for your deployment go to Management page where you can find and create keys in the 'Security section' -> API keys.

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-elastic

BATON_API_KEY=apiKey baton-elastic
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_API_KEY=apiKey ghcr.io/conductorone/baton-elastic:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-elastic/cmd/baton-elastic@main

BATON_API_KEY=apiKey baton-elastic
baton resources
```

# Data Model

`baton-elastic` will pull down information about the following Elastic resources:

By default:
- Users (if you want to sync only users of specific organization, provide the `--organization-id` flag, otherwise it syncs all users)
- Organizations

Optional: 
- Deployment roles
- Deployment users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-elastic` Command Line Usage

```
baton-elastic

Usage:
  baton-elastic [flags]
  baton-elastic [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-key string               Elastic API key used to communicate with Elastic cloud API. ($BATON_API_KEY)
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --deployment-api-key string    API key of your elasticsearch deployment. ($BATON_DEPLOYMENT_API_KEY)
      --deployment-endpoint string   Elasticsearch endpoint used to sync deployment resources. ($BATON_DEPLOYMENT_ENDPOINT)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-elastic
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --organization-id string       Optional. Provide your Elastic organization ID if you want to sync members of a single organization. ($BATON_ORGANIZATION_ID)
  -p, --provisioning                 This must be set in order for provisioning actions to be enabled. ($BATON_PROVISIONING)
  -v, --version                      version for baton-elastic

Use "baton-elastic [command] --help" for more information about a command.
```
