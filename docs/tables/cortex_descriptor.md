# Cortex Entity Table

This table calls the List entity descriptors API to get the data about each
entity descriptor (yaml definition). To see information about the entity from
all sources use the `entity` table.

## Examples

### Get information about a single entity descriptor

```sql
select
  tag,
  title,
  description,
  type,
  parents,
  groups,
  team,
  owners,
  slack,
  links,
  metadata
from
  cortex_descriptor
where
  tag = 'service1';
```
