-- +goose Up
-- +goose StatementBegin
ALTER TABLE system_preferences
ADD is_ai_enable BIT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE system_preferences
DROP COLUMN is_ai_enable;
-- +goose StatementEnd

