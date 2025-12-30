-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
ADD tenable_solution NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vulnerabilities
DROP COLUMN tenable_solution;
-- +goose StatementEnd
