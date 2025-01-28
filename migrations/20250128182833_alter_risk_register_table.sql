-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_register
ALTER COLUMN reference_number VARCHAR(45) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_register
ALTER COLUMN reference_number VARCHAR(45) NOT NULL;
-- +goose StatementEnd