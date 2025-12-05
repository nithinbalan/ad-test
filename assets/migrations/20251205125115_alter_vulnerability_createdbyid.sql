-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN created_by_id INT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN created_by_id INT NOT NULL;
-- +goose StatementEnd