-- +goose Up
-- +goose StatementBegin
-- Add the new related_task_ids column as NVARCHAR(MAX) for JSON compatibility
ALTER TABLE tasks
ADD related_task_ids NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the related_task_ids column
ALTER TABLE tasks
DROP COLUMN related_task_ids;
-- +goose StatementEnd