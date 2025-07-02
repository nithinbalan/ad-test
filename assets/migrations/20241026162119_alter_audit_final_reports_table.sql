-- +goose Up
-- +goose StatementBegin
-- Add the submitted_by_id column with a foreign key to users table
ALTER TABLE audit_final_reports
ADD submitted_by_id INT NULL;

-- Add the submitted_at column
ALTER TABLE audit_final_reports
ADD submitted_at DATETIME NULL;

-- Add the next_review_level column with a default value of 0
ALTER TABLE audit_final_reports
ADD next_review_level INT DEFAULT 0 NULL;

-- Add the foreign key constraint for submitted_by_id to users table
ALTER TABLE audit_final_reports
ADD CONSTRAINT FK_audit_final_reports_submitted_by
FOREIGN KEY (submitted_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for submitted_by_id
ALTER TABLE audit_final_reports
DROP CONSTRAINT FK_audit_final_reports_submitted_by;

-- Drop the submitted_by_id, submitted_at, and next_review_level columns
ALTER TABLE audit_final_reports
DROP COLUMN submitted_by_id;

ALTER TABLE audit_final_reports
DROP COLUMN submitted_at;

ALTER TABLE audit_final_reports
DROP COLUMN next_review_level;
-- +goose StatementEnd
