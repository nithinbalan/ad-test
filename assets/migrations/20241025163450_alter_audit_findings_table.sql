-- +goose Up
-- +goose StatementBegin
-- Add the department_response column
ALTER TABLE audit_findings
ADD department_response NVARCHAR(MAX) NULL;

-- Add the department_response_by_id column with a foreign key to users table
ALTER TABLE audit_findings
ADD department_response_by_id INT NULL;

-- Add the department_response_at column
ALTER TABLE audit_findings
ADD department_response_at DATETIME NULL;

-- Add the foreign key constraint for department_response_by_id to users
ALTER TABLE audit_findings
ADD CONSTRAINT FK_audit_findings_department_response_by
FOREIGN KEY (department_response_by_id) REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint first
ALTER TABLE audit_findings
DROP CONSTRAINT FK_audit_findings_department_response_by;

-- Drop the department_response, department_response_by_id, and department_response_at columns
ALTER TABLE audit_findings
DROP COLUMN department_response;

ALTER TABLE audit_findings
DROP COLUMN department_response_by_id;

ALTER TABLE audit_findings
DROP COLUMN department_response_at;
-- +goose StatementEnd
