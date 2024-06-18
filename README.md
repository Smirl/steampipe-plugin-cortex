<!-- Ensure this matches docs/index.md -->
# Cortex + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly
query cloud APIs using SQL.

[Cortex](https://cortex.io/) is the Internal Developer Portal eliminating
“developer tax” with paved paths to production. Create, catalog, score, and
drive action to continuously improve software.

For example:

```sql
select
    "name",
    "tag",
    "description",
    "type",
    "parents",
    "groups",
    "metadata",
    "last_updated",
    "links",
    "archived",
    "repository",
    "slack_channels"
from
    cortex_entity
limit 10
```

## Documentation

- **[Table definitions & examples →](/plugins/smirl/cortex/tables)**

## Get started

### Install

Download and install the latest AWS plugin:

```bash
steampipe plugin install smirl/cortex
```

### Credentials

You will need a Cortex API Token to authenticate with the API.

https://docs.cortex.io/docs/api/cortex-api

### Configuration

The configuration for the plugin is simple and requires just the API Token and
optional base url if you use self hosted cortex.

To connection to different instances, simply use a token from that other hosted
instance.

Environment variables can be used to override these configuration options.

```hcl
connection "cortex" {
    plugin    = "cortex"

    # API key from cortex.io for your instance
    # If the environment variable CORTEX_API_KEY is defined it will be overriden
    # api_key = "REPLACE_WITH_YOUR_CORTEX_API_KEY"

    # The BASE URL of your self hosted instance
    # If the environment variable CORTEX_BASE_URL is defined it will be overriden
    # base_url = "https://app.cortex.mycompany.com"
}
```