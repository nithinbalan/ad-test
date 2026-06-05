-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_level_types (    type,    title_en,    title_ar,    status,    module_identifier)
 VALUES (   'custom-roles',    'Custom Roles',   'أدوار مخصصة',    'active',    'organization');
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_level_types
WHERE type = 'custom-roles'
  AND module_identifier = 'organization';
-- +goose StatementEnd