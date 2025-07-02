-- +goose Up
-- +goose StatementBegin

-- Make auditable_item_id nullable
ALTER TABLE audit_findings
ALTER COLUMN auditable_item_id INT NULL;

-- Add new column audit_checklist_id
ALTER TABLE audit_findings
ADD audit_checklist_id INT NULL;

-- Add foreign key constraint for audit_checklist_id
ALTER TABLE audit_findings
ADD CONSTRAINT FK_audit_findings_audit_checklists
FOREIGN KEY (audit_checklist_id) REFERENCES audit_checklists(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign key constraint for audit_checklist_id
ALTER TABLE audit_findings
DROP CONSTRAINT FK_audit_findings_audit_checklists;

-- Drop the audit_checklist_id column
ALTER TABLE audit_findings
DROP COLUMN audit_checklist_id;

-- Revert auditable_item_id to NOT NULL if necessary
ALTER TABLE audit_findings
ALTER COLUMN auditable_item_id INT NOT NULL;

-- +goose StatementEnd
