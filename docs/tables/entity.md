# Cortex Entity Table

This table calls the "List entities" API to get the data about each entity. To
see the descriptor (yaml definition) of an entity use the `descriptor` table.


|Field|Description|
|----|-----|
|`name`|Pretty name of the entity.|
|`tag`|The x-cortex-tag of the entity.|
|`description`|Description.|
|`type`|Entity Type.|
|`parents`|Parents of the entity.|
|`groups`|Groups, kind of like tags.|
|`metadata`|Raw custom metadata|
|`last_updated`|Last updated time.|
|`links`|List of links|
|`archived`|Is archived.|
|`repository`|Git repo full name|
|`slack_channels`|List of string slack channels|

## Examples

### Get information about a single entity

```sql
select * from cortex_entity where tag = 'my-service'
```

### Count of all domains

```sql
select count(*) from cortex_entity where type = 'domain'
```

### Extract custom metadata into a column

```sql
select tag, metadata -> 'my_key' as my_key from cortex_entity limit 10
