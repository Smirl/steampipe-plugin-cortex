# Cortex Team Table

This table calls the List team API to get the data about each team. 


|Field|Description|
|----|-----|
|`name`|Pretty name of the team.|
|`tag`|The teamTag of the team.|
|`metadata`|Raw custom metadata|
|`links`|List of links|
|`archived`|Is archived.|
|`slack_channels`|List of string slack channels|
|`members`|List of team members|


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
from cortex_team
where tag = 'my-team'
```

### List teams for each member and the count of teams they are members of

```sql
select
    email,
    array_agg(tag) as teams,
    count(tag) as team_count
from (
    select
        tag,
        lower(jsonb_array_elements(members) ->> 'Email') as email
    from
        cortex_team
) as t
group by t.email
order by squad_count desc
```
