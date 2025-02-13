# Cortex Entity Table

This table calls the "List entities" API to get the data about each entity. To
see the descriptor (yaml definition) of an entity use the `descriptor` table.

## Examples

### Get information about a single entity

```sql
select
  *
from
  cortex_entity
where
  tag = 'service1'
```

```
+----------+----------+-----------------------------+---------+----------------------+-----------------+-----------------+---------------------------+--------------------------+----------+--------------+----------------+----------------+-------------------+
| name     | tag      | description                 | type    | parents              | groups          | metadata        | last_updated              | links                    | archived | repository   | slack_channels | owner_teams    | owner_individuals |
+----------+----------+-----------------------------+---------+----------------------+-----------------+-----------------+---------------------------+--------------------------+----------+--------------+----------------+----------------+-------------------+
| service1 | service1 | Test service for kubernetes | service | ["service1-project"] | ["third-party"] | {"key":"value"} | 2024-09-06T18:26:10+08:00 | ["https://example.com/"] | false    | k8s/service1 | <null>         | ["my-squad"]   | <null>            |
+----------+----------+-----------------------------+---------+----------------------+-----------------+-----------------+---------------------------+--------------------------+----------+--------------+----------------+----------------+-------------------+
```

### Count of all domains

```sql
select
  count(*)
from
  cortex_entity
where
  type = 'domain'
```

```
+-------+
| count |
+-------+
| 1856  |
+-------+
```

### Extract custom metadata into a column

```sql
select
  tag,
  metadata -> 'my_key' as my_key
from
  cortex_entity limit 10
```

```
+---------------------------+---------+
| tag                       | my_key  |
+---------------------------+---------+
| alpha-service             | value1  |
| beta-service              | value2  |
| gamma-service             | value3  |
| delta-service             | value4  |
| epsilon-service           | value5  |
| zeta-service              | value6  |
| eta-service               | value7  |
| theta-service             | value8  |
| iota-service              | value9  |
| kappa-service             | value10 |
+---------------------------+---------+
```
