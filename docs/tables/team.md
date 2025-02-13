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
  tag = 'my-team'
```

```
+-------------+-------------+----------------------------------------------------------+-------------------------------+----------+----------------+-----------------------------------------------------------------------------------------------+
| name        | tag         | metadata                                                 | links                         | archived | slack_channels | members                                                                                       |
+-------------+-------------+----------------------------------------------------------+-------------------------------+----------+----------------+-----------------------------------------------------------------------------------------------+
| Alpha Squad | alpha-squad | {"description":null,"name":"Alpha Squad","summary":null} | ["https://example.com/alpha"] | false    | <null>         | [{"Email":"joe.blogs@example.com","Name":"Joe Blogs","NotificationsEnabled":false,"Role":""}] |
+-------------+-------------+----------------------------------------------------------+-------------------------------+----------+----------------+-----------------------------------------------------------------------------------------------+
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
  team_count desc limit 10
```

```
+---------------------------+-----------------------------------------------------------------------------------+------------+
| email                     | teams                                                                             | team_count |
+---------------------------+-----------------------------------------------------------------------------------+------------+
| john.doe@example.com      | alpha-squad,beta-squad,gamma-squad,delta-squad,epsilon-squad,zeta-squad,eta-squad | 7          |
| jane.smith@example.com    | theta-squad,iota-squad,kappa-squad,lambda-squad,mu-squad,nu-squad                 | 6          |
| alice.jones@example.com   | xi-squad,omicron-squad,pi-squad,rho-squad,sigma-squad,tau-squad                   | 6          |
| bob.brown@example.com     | upsilon-squad,phi-squad,chi-squad,psi-squad,omega-squad                           | 5          |
| charlie.davis@example.com | alpha-squad,beta-squad,gamma-squad,delta-squad                                    | 4          |
| diana.evans@example.com   | epsilon-squad,zeta-squad,eta-squad,theta-squad                                    | 4          |
| emily.frank@example.com   | iota-squad,kappa-squad,lambda-squad                                               | 3          |
| frank.green@example.com   | mu-squad,nu-squad,xi-squad                                                        | 3          |
| grace.hill@example.com    | omicron-squad,pi-squad,rho-squad                                                  | 3          |
| hank.ivan@example.com     | sigma-squad,tau-squad,upsilon-squad                                               | 3          |
+---------------------------+-----------------------------------------------------------------------------------+------------+
```
