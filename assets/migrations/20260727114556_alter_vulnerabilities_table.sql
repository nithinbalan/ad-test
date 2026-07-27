-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN cve NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN cve NVARCHAR(100) NULL;
-- +goose StatementEnd
