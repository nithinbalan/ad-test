-- +goose Up
-- +goose StatementBegin
UPDATE integrations
SET 
    status = 'active',
    api_base_url = 'https://sc.tenalab.online'
WHERE type = 'tenable';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE integrations
SET 
    status = 'inactive',
    api_base_url = NULL
WHERE type = 'tenable';
-- +goose StatementEnd
