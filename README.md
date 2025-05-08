<!-- Ensure this matches docs/index.md -->
# Cortex + Steampipe

[![Coverage Status](https://coveralls.io/repos/github/Smirl/steampipe-plugin-cortex/badge.svg?branch=main)](https://coveralls.io/github/Smirl/steampipe-plugin-cortex?branch=main)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build](https://github.com/smirl/steampipe-plugin-cortex/actions/workflows/main.yaml/badge.svg)](https://github.com/Smirl/steampipe-plugin-cortex/actions/workflows/main.yaml)

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly
query cloud APIs using SQL.

[Cortex](https://cortex.io/) is the Internal Developer Portal eliminating
“developer tax” with paved paths to production. Create, catalog, score, and
drive action to continuously improve software.

For example:

```sql
select 
  tag,
  repository,
  owner_teams
from 
  cortex_entity 
where
  type = 'service'
limit 
  10;
```

## Documentation

- **[Table definitions & examples →](./docs/tables/)**

## Get started

### Install

Download and install the latest:

```bash
steampipe plugin install smirl/cortex
```

### Credentials

You will need a Cortex API Token to authenticate with the API.

https://docs.cortex.io/api/rest

### Configuration

The configuration for the plugin is simple and requires just the API Token and
optional base url if you use self hosted cortex.

To connection to different instances, simply use a token from that other hosted
instance.

Environment variables can be used to override these configuration options.

```hcl
connection "cortex" {
    plugin    = "smirl/cortex"

    # API key from cortex.io for your instance
    # If the environment variable CORTEX_API_KEY is defined it will be overriden
    # api_key = "REPLACE_WITH_YOUR_CORTEX_API_KEY"

    # The BASE URL of your self hosted instance
    # If the environment variable CORTEX_BASE_URL is defined it will be overriden
    # base_url = "https://app.cortex.mycompany.com"
}
```

## Get Involved

Open source: https://github.com/smirl/steampipe-plugin-cortex
