-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_custom_roles (
    title_en,
    title_ar,
    type,
    module_identifier
)
VALUES (
    'Risk Responsible User',
    'المستخدم المسؤول عن المخاطر',
    'risk_responsible_user',
    'asset'
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_custom_roles
WHERE type = 'risk_responsible_user'
  AND module_identifier = 'asset';
-- +goose StatementEnd
