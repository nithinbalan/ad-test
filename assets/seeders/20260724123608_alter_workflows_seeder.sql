-- +goose Up
-- +goose StatementBegin

UPDATE workflows
SET type = 'policy-workflow'
WHERE title_en = 'Policy Workflow';

UPDATE workflows
SET type = 'procedure-workflow'
WHERE title_en = 'Procedure Workflow';

UPDATE workflows
SET type = 'risk-workflow'
WHERE title_en = 'Risk Workflow';

UPDATE workflows
SET type = 'assessment-workflow'
WHERE title_en = 'Assessment Workflow';

UPDATE workflows
SET type = 'committee-management-workflow'
WHERE title_en = 'Committee Management Workflow';

UPDATE workflows
SET type = 'kri-register-workflow'
WHERE title_en = 'KRI Register Workflow';

UPDATE workflows
SET type = 'risk-response-workflow'
WHERE title_en = 'Risk Response Workflow';

UPDATE workflows
SET type = 'risk-register-update-plan-workflow'
WHERE title_en = 'Risk Register Update Plan Workflow';

UPDATE workflows
SET type = 'compliance-register-workflow'
WHERE title_en = 'Compliance Register Workflow';

UPDATE workflows
SET type = 'bia-workflow'
WHERE title_en = 'BIA Workflow';

UPDATE workflows
SET type = 'risk-acceptance-workflow'
WHERE title_en = 'Risk Acceptance Workflow';

UPDATE workflows
SET type = 'vulnerability-closure-workflow'
WHERE title_en = 'Vulnerability Closure Workflow';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE workflows
SET type = NULL
WHERE title_en IN (
    'Policy Workflow',
    'Procedure Workflow',
    'Risk Workflow',
    'Assessment Workflow',
    'Committee Management Workflow',
    'KRI Register Workflow',
    'Risk Response Workflow',
    'Risk Register Update Plan Workflow',
    'Compliance Register Workflow',
    'BIA Workflow',
    'Risk Acceptance Workflow',
    'Vulnerability Closure Workflow'
);

-- +goose StatementEnd