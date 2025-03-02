-- +goose Up
-- +goose StatementBegin
ALTER TABLE integration_credentials 
ALTER COLUMN access_token NVARCHAR(MAX);

ALTER TABLE integration_credentials 
ALTER COLUMN refresh_token NVARCHAR(MAX);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE integration_credentials 
ALTER COLUMN access_token NVARCHAR(255); 
ALTER TABLE integration_credentials 
ALTER COLUMN refresh_token NVARCHAR(255); 
-- +goose StatementEnd

