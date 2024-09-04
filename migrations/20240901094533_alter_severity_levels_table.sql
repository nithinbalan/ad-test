-- +goose Up
-- +goose StatementBegin
ALTER TABLE severity_levels ADD weight INT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE severity_levels DROP COLUMN weight;
-- +goose StatementEnd
