-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_statuses
ADD is_active BIT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_statuses
DROP COLUMN is_active;
-- +goose StatementEnd
