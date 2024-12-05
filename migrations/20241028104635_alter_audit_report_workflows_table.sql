-- +goose Up
-- +goose StatementBegin
-- Add the workflow_level_type_id column with a foreign key to workflow_level_types table
ALTER TABLE audit_report_workflows
ADD workflow_level_type_id INT NULL;

-- Add the foreign key constraint for workflow_level_type_id to workflow_level_types
ALTER TABLE audit_report_workflows
ADD CONSTRAINT FK_audit_report_workflows_level_types
FOREIGN KEY (workflow_level_type_id) REFERENCES workflow_level_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for workflow_level_type_id
ALTER TABLE audit_report_workflows
DROP CONSTRAINT FK_audit_report_workflows_level_types;

ALTER TABLE audit_report_workflows
DROP COLUMN workflow_level_type_id;
-- +goose StatementEnd