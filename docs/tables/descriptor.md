# Cortex Entity Table

This table calls the List entity descriptors API to get the data about each
entity descriptor (yaml definition). To see information about the entity from
all sources use the `entity` table.


|Field|Description|
|----|-----|
|`tag`|The x-cortex-tag of the entity|
|`title`|Title|
|`description`|Description|
|`type`|Entity Type|
|`parents`|Parent tags.|
|`groups`|Groups, kind of like tags|
|`team`|Raw team|
|`owners`|Raw owner|
|`slack`|Raw slack|
|`links`|List of links|
|`metadata`|Raw custom metadata|
|`repository`|Git repo full name|
|`victorops`|Victorops team slug|
|`jira`|List of jira projects|

## Examples

### Get information about a single entity descriptor

```sql
select * from cortex_descriptor where tag = 'my-service'
```
