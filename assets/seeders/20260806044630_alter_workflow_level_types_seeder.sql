-- +goose Up
-- +goose StatementBegin
INSERT INTO workflow_level_types (type, title_en, title_ar, status, module_identifier)
VALUES ('user-and-reporting-user', 'User and Reporting User', 'المستخدم والمستخدم المُبلغ عنه', 'active', 'organization,compliance');
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DELETE FROM workflow_level_types
WHERE type = 'user-and-reporting-user'
  AND module_identifier = 'organization,compliance';
-- +goose StatementEnd