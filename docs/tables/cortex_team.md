# Cortex Team Table

This table calls the List team API to get the data about each team. 

## Examples

### Get information about a team

```sql
select
  name,
  tag,
  metadata,
  links,
  archived,
  slack_channels,
  members 
from
  cortex_team 
where
  tag = 'my-team';
```

### List teams for each member and the count of teams they are members of

```sql
select
  email,
  array_agg(tag) as teams,
  count(tag) as team_count 
from
  (
    select
      tag,
      lower(jsonb_array_elements(members) ->> 'Email') as email 
    from
      cortex_team 
  )
  as t 
group by
  t.email 
order by
  team_count desc 
limit 
  10;
```
