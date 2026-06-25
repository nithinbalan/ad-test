-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET title_en = 'Corrective Action',
    title_ar = N'الإجراء التصحيحي',
    route = '/compliance/corrective-actions'
WHERE identifier = 'compliance_action_plan';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET title_en = 'Action Plans',
    title_ar = N'خطط العمل',
    route = '/compliance/action-plans'
WHERE identifier = 'compliance_action_plan';
-- +goose StatementEnd
