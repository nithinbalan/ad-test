-- +goose Up
-- +goose StatementBegin
-- Add the new audit_document_request_id column as INT
ALTER TABLE tasks
ADD audit_document_request_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the audit_document_request_id column
ALTER TABLE tasks
DROP COLUMN audit_document_request_id;
-- +goose StatementEnd