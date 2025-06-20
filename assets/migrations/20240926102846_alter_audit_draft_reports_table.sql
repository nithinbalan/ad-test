-- +goose Up
-- +goose StatementBegin

-- Dropping existing fields
ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_users_3;

ALTER TABLE audit_draft_reports
DROP COLUMN submitted_by_id;

ALTER TABLE audit_draft_reports
DROP COLUMN submitted_at;

-- Adding new fields for internal review and auditee review
ALTER TABLE audit_draft_reports
ADD internal_review_submitted_by_id INT NULL;

ALTER TABLE audit_draft_reports
ADD internal_review_submitted_at DATETIME NULL;

ALTER TABLE audit_draft_reports
ADD auditee_review_submitted_by_id INT NULL;

ALTER TABLE audit_draft_reports
ADD auditee_review_submitted_at DATETIME NULL;

-- Adding foreign keys for the new fields
ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_internal_review_submitted_by_id FOREIGN KEY (internal_review_submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_auditee_review_submitted_by_id FOREIGN KEY (auditee_review_submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- Dropping foreign key constraints for the new fields
ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_internal_review_submitted_by_id;

ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_auditee_review_submitted_by_id;

-- Dropping the new fields
ALTER TABLE audit_draft_reports
DROP COLUMN internal_review_submitted_by_id;

ALTER TABLE audit_draft_reports
DROP COLUMN internal_review_submitted_at;

ALTER TABLE audit_draft_reports
DROP COLUMN auditee_review_submitted_by_id;

ALTER TABLE audit_draft_reports
DROP COLUMN auditee_review_submitted_at;

-- Re-adding the original fields
ALTER TABLE audit_draft_reports
ADD submitted_by_id INT NULL;

ALTER TABLE audit_draft_reports
ADD submitted_at DATETIME NULL;

-- Re-adding the foreign key constraint for submitted_by_id
ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_users_3 FOREIGN KEY (submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
