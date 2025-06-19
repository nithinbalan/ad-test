-- +goose Up
-- +goose StatementBegin

-- Adding new fields for completed details
ALTER TABLE audit_draft_reports
ADD completed_at DATETIME NULL;

ALTER TABLE audit_draft_reports
ADD completed_by_id INT NULL;

-- Adding foreign key constraint for completed_by_id
ALTER TABLE audit_draft_reports
ADD CONSTRAINT FK_audit_draft_reports_completed_by_id FOREIGN KEY (completed_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Dropping foreign key constraint for completed_by_id
ALTER TABLE audit_draft_reports
DROP CONSTRAINT FK_audit_draft_reports_completed_by_id;

-- Dropping the new fields
ALTER TABLE audit_draft_reports
DROP COLUMN completed_at;

ALTER TABLE audit_draft_reports
DROP COLUMN completed_by_id;

-- +goose StatementEnd
