# Cortex Entity Table

This table calls the List entity descriptors API to get the data about each
entity descriptor (yaml definition). To see information about the entity from
all sources use the `entity` table.

## Examples

### Get information about a single entity descriptor

```sql
select
  *
from
  cortex_descriptor
where
  tag = 'service1'
```

```
+----------+----------+-----------------------------+---------+---------------+------------+--------+--------+-------------------------------------------------------------------------------------------------------+--------+--------------------------------------------------------------+
| tag      | title    | description                 | type    | parents       | groups     | team   | owners | slack                                                                                                 | links  | metadata         | repository   | victorops       | jira     | 
+----------+----------+-----------------------------+---------+---------------+------------+--------+--------+-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| service1 | Service1 | Test service for kubernetes | service | ["my-domain"] | ["groupa"] | <null> | <null> | {"Channels":[{"Description":"Slack Channel","Name":"service1-support","NotificationsEnabled":false}]} | <null> | {"key": "value"} | org/service1 | victorops-team1 | ["JIRA"] |
+----------+----------+-----------------------------+---------+---------------+------------+--------+--------+-------------------------------------------------------------------------------------------------------+--------+--------------------------------------------------------------+
```
