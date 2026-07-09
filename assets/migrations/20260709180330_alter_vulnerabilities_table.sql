-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
    ADD justification NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
    DROP COLUMN justification;
-- +goose StatementEnd