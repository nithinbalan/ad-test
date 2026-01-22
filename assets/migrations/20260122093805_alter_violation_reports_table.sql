-- +goose Up
-- +goose StatementBegin
ALTER TABLE violation_reports
ADD email_of_reporter NVARCHAR(MAX) NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE violation_reports
DROP COLUMN email_of_reporter;
-- +goose StatementEnd
