# Scorecard Scores Table

## Description
Description of the Scorecard Scores table.

## SQL Examples

### Query scores for a specific scorecard

```sql
select
  scorecard_tag,
  scorecard_name,
  service_tag,
  service_name,
  service_groups,
  last_evaluated,
  rule_identifier,
  rule_title,
  rule_description,
  rule_expression,
  rule_effective_from,
  rule_level_name,
  rule_level_number,
  rule_weight,
  rule_score,
  rule_pass 
from
  cortex_scorecard_score
where
  scorecard_tag = 'my-scorecard'
limit 
  10;
```

### Query scores where the rule passed

```sql
select
  service_tag,
  rule_expression
from
  cortex_scorecard_scores
where
  scorecard_tag = 'my-scorecard' 
  and rule_pass = true 
order by
  rule_level_number
limit 
  10;
```
