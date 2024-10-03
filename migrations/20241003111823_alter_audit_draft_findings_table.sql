-- +goose Up
-- +goose StatementBegin

-- Adding new fields for accepted and challenged states
ALTER TABLE audit_draft_findings
ADD accepted_by_id INT NULL;

ALTER TABLE audit_draft_findings
ADD accepted_at DATETIME NULL;

ALTER TABLE audit_draft_findings
ADD challenged_by_id INT NULL;

ALTER TABLE audit_draft_findings
ADD challenged_at DATETIME NULL;

-- Adding foreign keys for the new fields
ALTER TABLE audit_draft_findings
ADD CONSTRAINT FK_audit_draft_findings_accepted_by_id FOREIGN KEY (accepted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE audit_draft_findings
ADD CONSTRAINT FK_audit_draft_findings_challenged_by_id FOREIGN KEY (challenged_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- Dropping foreign key constraints for the new fields
ALTER TABLE audit_draft_findings
DROP CONSTRAINT FK_audit_draft_findings_accepted_by_id;

ALTER TABLE audit_draft_findings
DROP CONSTRAINT FK_audit_draft_findings_challenged_by_id;

-- Dropping the new fields
ALTER TABLE audit_draft_findings
DROP COLUMN accepted_by_id;

ALTER TABLE audit_draft_findings
DROP COLUMN accepted_at;

ALTER TABLE audit_draft_findings
DROP COLUMN challenged_by_id;

ALTER TABLE audit_draft_findings
DROP COLUMN challenged_at;

-- +goose StatementEnd
