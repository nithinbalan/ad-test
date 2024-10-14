-- +goose Up
-- +goose StatementBegin
-- Add the new document_workflow_id column as INT 
ALTER TABLE tasks
ADD document_workflow_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the document_workflow_id column
ALTER TABLE tasks
DROP COLUMN document_workflow_id;
-- +goose StatementEnd