-- +goose Up
-- +goose StatementBegin
-- Add the role_id column with a foreign key to roles table
ALTER TABLE audit_report_workflow_participants
ADD role_id INT NULL;

-- Add updated_by_id and updated_at columns
ALTER TABLE audit_report_workflow_participants
ADD updated_by_id INT NULL;

ALTER TABLE audit_report_workflow_participants
ADD updated_at DATETIME NULL;

-- Add the foreign key constraint for role_id to roles table
ALTER TABLE audit_report_workflow_participants
ADD CONSTRAINT FK_audit_report_workflow_participants_roles
FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE audit_report_workflow_participants
ADD CONSTRAINT FK_audit_report_workflow_participants_updated_by_users
FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for role_id
ALTER TABLE audit_report_workflow_participants
DROP CONSTRAINT FK_audit_report_workflow_participants_roles;

-- Drop the role_id
ALTER TABLE audit_report_workflow_participants
DROP COLUMN role_id;

-- Remove the foreign key constraint for updated_by_id
ALTER TABLE audit_report_workflow_participants
DROP CONSTRAINT FK_audit_report_workflow_participants_updated_by_users;

-- Drop the updated_by_id
ALTER TABLE audit_report_workflow_participants
DROP COLUMN updated_by_id;

-- Drop the updated_at
ALTER TABLE audit_report_workflow_participants
DROP COLUMN updated_at;

-- +goose StatementEnd
