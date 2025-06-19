-- +goose Up
-- +goose StatementBegin
ALTER TABLE system_preferences
ADD license_key NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE system_preferences
DROP COLUMN license_key;
-- +goose StatementEnd
