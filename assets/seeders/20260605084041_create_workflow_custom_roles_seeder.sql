-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_custom_roles (
    title_en,
    title_ar,
    type,
    module_identifier
)
VALUES (
    'Department Head',
    'رئيس القسم',
    'department_head',
    'organization'
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_custom_roles
WHERE type = 'department_head'
  AND module_identifier = 'organization';
-- +goose StatementEnd