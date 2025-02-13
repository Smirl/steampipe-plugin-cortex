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
limit 10
```

```
+-------------------+-------------------+-------------------+-------------------+---------------------+----------------------------+--------------------------------------+------------------------------------------------+-------------------------------------------------+-----------------+-------------------+-------------+------------+-----------+
| scorecard_tag     | scorecard_name    | service_tag       | service_name      | service_groups      | last_evaluated             | rule_identifier                      | rule_title                                     | rule_expression                                 | rule_level_name | rule_level_number | rule_weight | rule_score | rule_pass |
+-------------------+-------------------+-------------------+-------------------+---------------------+----------------------------+--------------------------------------+------------------------------------------------+-------------------------------------------------+-----------------+-------------------+-------------+------------+-----------+
| service-readiness | Cortex Onboarding | service-1         | service-1         | ["microservice"]    | 2025-02-12T21:09:18.246024 | c536b79c-f35e-3e47-833f-194df4860f77 | Updated by Github App                          | custom("updated-by") != null                    | Bronze          | 1                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-1         | service-1         | ["microservice"]    | 2025-02-12T21:09:18.246024 | 92ab2234-ffb1-371b-9488-8ebe3485e503 | Has owners                                     | ownership.allOwners().length > 0                | Bronze          | 1                 | 5           | 5          | true      |
| service-readiness | Cortex Onboarding | service-2         | service-2         | ["component"]       | 2025-02-12T14:59:52.951055 | 92ab2234-ffb1-371b-9488-8ebe3485e503 | Has owners                                     | ownership.allOwners().length > 0                | Bronze          | 1                 | 5           | 5          | true      |
| service-readiness | Cortex Onboarding | service-3         | service-3         | ["third-party"]     | 2025-02-12T22:15:00.232289 | e423f153-809f-3279-ac30-4ee48aa8ca95 | soxScope field set in catalog.yml              | custom("sox-scope") != "undefined"              | Gold            | 3                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-1         | service-1         | ["microservice"]    | 2025-02-12T21:09:18.246024 | e423f153-809f-3279-ac30-4ee48aa8ca95 | soxScope field set in catalog.yml              | custom("sox-scope") != "undefined"              | Gold            | 3                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-1         | service-1         | ["microservice"]    | 2025-02-12T21:09:18.246024 | 2a235b74-711b-396a-a8da-994e437b5ee8 | personalDataHandling field set in catalog.yml  | custom("personal-data-handling") != "undefined" | Gold            | 3                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-1         | service-1         | ["microservice"]    | 2025-02-12T21:09:18.246024 | 60e0fa9e-bc63-3f25-9a08-0e50b271494b | paymentProcessing field set in catalog.yml     | custom("payment-processing") != "undefined"     | Gold            | 3                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-2         | service-2         | ["component"]       | 2025-02-12T14:59:52.951055 | 945e4ff1-f775-3f3b-8b66-7199e1a595ab | Has description                                | entity.description() != null                    | Bronze          | 1                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-4         | service-4         | ["lambda-function"] | 2025-02-13T04:37:38.739684 | c536b79c-f35e-3e47-833f-194df4860f77 | Updated by Github App                          | custom("updated-by") != null                    | Bronze          | 1                 | 1           | 1          | true      |
| service-readiness | Cortex Onboarding | service-5         | service-5         | ["other"]           | 2025-02-12T21:41:45.857825 | 945e4ff1-f775-3f3b-8b66-7199e1a595ab | Has description                                | entity.description() != null                    | Bronze          | 1                 | 1           | 1          | true      |
+-------------------+-------------------+-------------------+-------------------+---------------------+----------------------------+--------------------------------------+------------------------------------------------+-------------------------------------------------+-----------------+-------------------+-------------+------------+-----------+
```

# Query scores where the rule passed

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
limit 10
```

```
+-------------------+------------------------------------------------------------+
| service_tag       | rule_expression                                            |
+-------------------+------------------------------------------------------------+
| service-1         | git != null                                                |
| service-1         | git.fileExists("README.md") or git.fileExists("readme.md") |
| service-1         | ownership.allOwners().length > 0                           |
| service-2         | custom("updated-by") != null                               |
| service-2         | entity.description() != null                               |
| service-2         | git != null                                                |
| service-2         | ownership.allOwners().length > 0                           |
| service-3         | custom("updated-by") != null                               |
| service-3         | entity.description() != null                               |
| service-3         | git != null                                                |
+-------------------+------------------------------------------------------------+
```
