-- +goose Up
-- +goose StatementBegin


ALTER TABLE workflow_level_types
DROP CONSTRAINT CHK_workflow_level_types_type;


ALTER TABLE workflow_level_types
ADD CONSTRAINT CHK_workflow_level_types_type
CHECK (type IN ('individual', 'multiple-users', 'role', 'hod', 'custom-roles'));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin


ALTER TABLE workflow_level_types
DROP CONSTRAINT CHK_workflow_level_types_type;


ALTER TABLE workflow_level_types
ADD CONSTRAINT CHK_workflow_level_types_type
CHECK (type IN ('individual', 'multiple-users', 'role', 'hod'));
-- +goose StatementEnd
