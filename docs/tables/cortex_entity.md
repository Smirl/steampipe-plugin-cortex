# Cortex Entity Table

This table calls the "List entities" API to get the data about each entity. To
see the descriptor (yaml definition) of an entity use the `descriptor` table.

## Examples

### Get information about a single entity

```sql
select
  name,
  tag,
  description,
  type,
  parents,
  groups,
  metadata
from
  cortex_entity
where
  tag = 'service1';
```

### Count of all domains

```sql
select
  count(*)
from
  cortex_entity
where
  type = 'domain';
```

### Extract custom metadata into a column

```sql
select
  tag,
  metadata -> 'my_key' as my_key
from
  cortex_entity 
limit 
  10;
```
