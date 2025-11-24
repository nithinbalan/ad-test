-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles
ADD type VARCHAR(50) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE roles
DROP COLUMN type;
-- +goose StatementEnd
