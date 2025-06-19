-- +goose Up
-- +goose StatementBegin
DECLARE @constraint_name NVARCHAR(200)
SELECT @constraint_name = name
FROM sys.check_constraints
WHERE parent_object_id = OBJECT_ID('grc_report_sections')
AND type = 'C';

IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE grc_report_sections DROP CONSTRAINT ' + @constraint_name)
END

ALTER TABLE grc_report_sections
ADD CONSTRAINT CK_grc_report_sections_type CHECK (type IN (
    'main-executive-summary',
    'governance-summary',
    'executive-summary-governance',
    'policy-approved-overview',
    'governance-tasks',
    'conclusion-governance',
    'executive-summary-compliance',
    'compliance-summary-self-assessment',
    'overall-compliance-status-by-department-self-assessment',
    'conclusion-recommendations-compliance',
    'integrated-grc-performance',
    'risk-management-summary',
    'executive-summary-risk',
    'risk-governance-context',
    'top-5-risks',
    'key-risk-indicators-kris-performance',
    'risk-control-matrix',
    'risk-monitoring-reporting-communication',
    'recommendations-action-plans',
    'conclusion-recommendations-risk',
    'main-conclusion',
    'sla-summary'
));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @constraint_name NVARCHAR(200)
SELECT @constraint_name = name
FROM sys.check_constraints
WHERE parent_object_id = OBJECT_ID('grc_report_sections')
AND type = 'C';

IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE grc_report_sections DROP CONSTRAINT ' + @constraint_name)
END

ALTER TABLE grc_report_sections
ADD CONSTRAINT CK_grc_report_sections_type CHECK (type IN (
    'main-executive-summary',
    'governance-summary',
    'executive-summary-governance',
    'policy-approved-overview',
    'governance-tasks',
    'conclusion-governance',
    'executive-summary-compliance',
    'compliance-summary-self-assessment',
    'overall-compliance-status-by-department-self-assessment',
    'conclusion-recommendations-compliance',
    'integrated-grc-performance',
    'risk-management-summary',
    'executive-summary-risk',
    'risk-governance-context',
    'top-5-risks',
    'key-risk-indicators-kris-performance',
    'risk-control-matrix',
    'risk-monitoring-reporting-communication',
    'recommendations-action-plans',
    'conclusion-recommendations-risk',
    'main-conclusion'
));
-- +goose StatementEnd
