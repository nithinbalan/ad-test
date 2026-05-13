-- +goose Up
-- +goose StatementBegin

ALTER TABLE modules ADD is_passive BIT NOT NULL CONSTRAINT DF_modules_is_passive DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE modules DROP CONSTRAINT DF_modules_is_passive;
ALTER TABLE modules DROP COLUMN is_passive;

-- +goose StatementEnd
