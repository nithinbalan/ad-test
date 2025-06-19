-- +goose Up
-- +goose StatementBegin
-- Add the published_by_id column with a foreign key to users table
ALTER TABLE audit_final_reports
ADD published_by_id INT NULL;

-- Add the published_at column
ALTER TABLE audit_final_reports
ADD published_at DATETIME NULL;

-- Add the comments column
ALTER TABLE audit_final_reports
ADD comments NVARCHAR(MAX) NULL;

-- Add the foreign key constraint for published_by_id to users table
ALTER TABLE audit_final_reports
ADD CONSTRAINT FK_audit_final_reports_published_by
FOREIGN KEY (published_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for published_by_id
ALTER TABLE audit_final_reports
DROP CONSTRAINT FK_audit_final_reports_published_by;

-- Drop the published_by_id, published_at, and comments columns
ALTER TABLE audit_final_reports
DROP COLUMN published_by_id;

ALTER TABLE audit_final_reports
DROP COLUMN published_at;

ALTER TABLE audit_final_reports
DROP COLUMN comments;
-- +goose StatementEnd
