-- +goose Up
-- +goose StatementBegin
ALTER TABLE audits
ADD audit_status_id INT NULL; -- Add audit_status_id column

-- Add foreign key constraint to reference audit_statuses table for audit_status_id
ALTER TABLE audits
ADD CONSTRAINT FK_audits_audit_statuses FOREIGN KEY (audit_status_id) REFERENCES audit_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraint
ALTER TABLE audits
DROP CONSTRAINT FK_audits_audit_statuses;

-- Drop audit_status_id column
ALTER TABLE audits
DROP COLUMN audit_status_id;
-- +goose StatementEnd
