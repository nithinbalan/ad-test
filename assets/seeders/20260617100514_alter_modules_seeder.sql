-- +goose Up
-- +goose StatementBegin

UPDATE modules
SET is_passive = 1
WHERE identifier IN (
    'risk_mitigation',
    'compliance_action_plan',
    'department_response',
    'incident_action_plan',
    'meeting_action_plan',
    'policies'
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE modules
SET is_passive = 0
WHERE identifier IN (
    'risk_mitigation',
    'compliance_action_plan',
    'department_response',
    'incident_action_plan',
    'meeting_action_plan',
    'policies'
);

-- +goose StatementEnd
