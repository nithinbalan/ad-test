# QA Guide: Compliance Report Scoring & New Fields

## 1. Overview

The **Compliance Report Details** API now surfaces the following information for every requirement inside a compliance report:

- Average compliance score
- Current and average compliance status
- Inherent and residual risk scores
- Inherent and residual risk rating IDs and their title/color/type
- Whether the assessment is risk-based or regulation-notification based
- Acknowledgement status, user, and comment

This document explains how each value is calculated and where it comes from.

## 2. Data Sources

The report is built from three main sources:

| Source | Model | Key fields used |
|---|---|---|
| Requirement definition | `ComplianceRequirement` | `InherentRiskScore`, `InherentRiskRatingID`, `ResidualRiskRatingID`, `ComplianceStatusID`, `TitleEn`, `TitleAr` |
| Assessment requirement | `ComplianceAssessmentRequirement` | `AverageResidualRiskScore`, `ComplianceStatusID` |
| Responsible unit answer | `ComplianceAssessmentRequirementResponsibleUnit` | `Score`, `ComplianceStatusID`, `IsNotApplicable`, `IsRiskBasedAssessment`, `IsRegulationNotificationAssessment`, `IsAcknowledged`, `AcknowledgeUserID`, `AcknowledgeComment` |

## 3. Scoring Calculation

### 3.1 Framework detection

Previously the code tried to detect the scoring framework from `ComplianceAssessmentFrameworkValue`. In practice this relationship is empty, so the report now defaults to a **percentage-based** calculation.

- If any responsible-unit record has a `Score`, the framework is treated as **percentage** with a maximum of `100`.
- If no scores exist, it still defaults to percentage so the API does not break.

### 3.2 Average score for a requirement

For every requirement, all responsible-unit answers are collected and the average is computed as follows:

```
sum of all applicable Score values
----------------------------------  x 100
 count of applicable answers x 100
```

- Only rows where `IsNotApplicable = false` are included.
- The result is rounded to two decimals.
- A perfect score is stored as `99.99` rather than `100`.

### 3.3 Average compliance status for a requirement

From the same set of responsible-unit answers, the **most frequent** `ComplianceStatusID` is chosen. If no status exists, the requirement’s own current compliance status is used as a fallback.

### 3.4 Department compliance percentage

For a department (responsible unit) report, the department percentage is the average of the requirement averages described above. The department’s compliance status is derived from that percentage using the master `ComplianceStatus` ranges.

## 4. Risk Scores & Ratings

### 4.1 Inherent values

| Field | Source | Description |
|---|---|---|
| `inherent_score` | `ComplianceRequirement.InherentRiskScore` | The inherent risk score recorded on the requirement itself. |
| `inherent_risk_rating_id` | `ComplianceRequirement.InherentRiskRatingID` | The linked risk rating ID from the requirement. |
| `inherent_risk_rating` | `risk_matrix_risk_ratings` table via FK | The title, Arabic title, color code, and type of the inherent rating. |

### 4.2 Residual values

| Field | Source | Description |
|---|---|---|
| `residual_score` | `ComplianceAssessmentRequirement.AverageResidualRiskScore` first; if empty, computed from responsible-unit `Score` values. | The remaining risk after assessment. |
| `residual_risk_rating_id` | `ComplianceRequirement.ResidualRiskRatingID` | The linked residual risk rating ID from the requirement. |
| `residual_risk_rating` | `risk_matrix_risk_ratings` table via FK | The title, Arabic title, color code, and type of the residual rating. |

## 5. Acknowledgement Fields

A requirement can be marked as acknowledged at the responsible-unit level.

| Field | Source | Description |
|---|---|---|
| `is_acknowledged` | Aggregated from responsible units (`IsAcknowledged`) | `true` if any responsible unit has acknowledged. |
| `acknowledge_user_id` | From the first acknowledged responsible unit | The user who performed the acknowledgement. |
| `acknowledge_comment` | From the first acknowledged responsible unit | The comment entered during acknowledgement. |
| `acknowledge_user` | User preload | Name and email of the acknowledging user. |

## 6. Assessment Type Flags

Two flags are now stored on each report requirement and department to show what kind of assessment generated the data.

| Field | Source | Description |
|---|---|---|
| `is_risk_based_assessment` | `ComplianceAssessmentRequirementResponsibleUnit.IsRiskBasedAssessment` | `true` when the assessment was created as a risk-based assessment. |
| `is_regulation_notification_assessment` | `ComplianceAssessmentRequirementResponsibleUnit.IsRegulationNotificationAssessment` | `true` when the assessment was created from a regulation notification. |

These flags are also present on the department level so the report can show, per department, whether it came from a risk-based or regulation-notification workflow.

## 7. What Changed in the Code

### 7.1 Scoring logic

- `inferFrameworkMetaFromAssessments` no longer relies on `ComplianceAssessmentFrameworkValue`. It now falls back to percentage scoring whenever a `Score` exists.
- `avgScoreAndStatusFromRUsFrameworkAware` uses `Score` and `ComplianceStatusID` directly from responsible-unit records and skips `IsNotApplicable` rows.
- `deptAverageForUnit` uses `IsNotApplicable` instead of the empty framework-value relationship.

### 7.2 New database columns

Migration `20250827000000_alter_compliance_report_requirements_add_risk_ratings.sql` adds:

- `compliance_report_requirements.inherent_risk_rating_id`
- `compliance_report_requirements.residual_risk_rating_id`

Both are nullable foreign keys to `risk_matrix_risk_ratings(id)`.

### 7.3 Preloads

The report-details repository now preloads:

- `ComplianceReportRequirements.InherentRiskRating`
- `ComplianceReportRequirements.ResidualRiskRating`
- Same for `ComplianceReportNonCompliantRequirements`

### 7.4 DTO additions

`ComplianceReportRequirement` DTO now includes:

- `inherent_risk_rating_id`
- `inherent_risk_rating` (object)
- `residual_risk_rating_id`
- `residual_risk_rating` (object)

## 8. How to Verify

1. Create or regenerate a compliance report.
2. Call the **Get Full Details** endpoint.
3. For each requirement confirm:
   - `average_score` is populated when responsible units have scores.
   - `average_compliance_status` matches the most common status.
   - `inherent_score`, `residual_score`, `inherent_risk_rating_id`, and `residual_risk_rating_id` are populated.
   - `inherent_risk_rating` and `residual_risk_rating` objects contain title and color.
   - `is_risk_based_assessment` / `is_regulation_notification_assessment` reflect the assessment type.
   - `is_acknowledged` and `acknowledge_user` reflect acknowledgement data.

## 9. Common Questions

**Q: Why was `ComplianceAssessmentFrameworkValue` removed from the calculation?**  
A: In the live dataset this relationship is always empty (`NULL`), so depending on it caused every average score to become `0`. The calculation now uses the `Score` column, which is populated.

**Q: What happens if a responsible unit is marked not applicable?**  
A: It is excluded from both the average score and the modal status calculation.

**Q: What happens if there is no residual risk rating ID?**  
A: The field remains `NULL`. The `residual_score` is still calculated from the assessment answers when possible.

**Q: Are these changes backward compatible?**  
A: Yes. All new columns are nullable, and the scoring change only affects the calculation path, not the database schema for existing data.
