-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ADD tenable_cve VARCHAR(100) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
DROP COLUMN tenable_cve;
-- +goose StatementEnd
