-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN description NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ALTER COLUMN description NVARCHAR(500) NULL;
-- +goose StatementEnd
