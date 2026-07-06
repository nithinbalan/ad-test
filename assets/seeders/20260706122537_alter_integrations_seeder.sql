-- +goose Up
-- +goose StatementBegin
UPDATE integrations
SET 
    status = 'active'
WHERE type = 'qualys';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE integrations
SET 
    status = 'inactive'
WHERE type = 'qualys';
-- +goose StatementEnd
