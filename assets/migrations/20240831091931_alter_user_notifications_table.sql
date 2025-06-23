-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_notifications
ADD entity NVARCHAR(255) NULL,
    action NVARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_notifications
DROP COLUMN entity,
DROP COLUMN action;
-- +goose StatementEnd
