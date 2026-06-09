-- +goose Up
-- +goose StatementBegin
UPDATE workflow_statuses
SET title_en = 'Reopened'
WHERE type = 'reopen';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE workflow_statuses
SET title_en = 'ReOpened'
WHERE type = 'reopen';
-- +goose StatementEnd
