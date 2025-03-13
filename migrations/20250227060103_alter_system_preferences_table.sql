-- +goose Up
-- +goose StatementBegin
ALTER TABLE system_preferences
ADD license_key_expiry_date DATE NULL;

ALTER TABLE system_preferences
ADD total_licensed_user_count INT NULL;

ALTER TABLE system_preferences
ADD balance_licensed_user_count INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE system_preferences
DROP COLUMN license_key_expiry_date;

ALTER TABLE system_preferences
DROP COLUMN total_licensed_user_count;

ALTER TABLE system_preferences
DROP COLUMN balance_licensed_user_count;
-- +goose StatementEnd
