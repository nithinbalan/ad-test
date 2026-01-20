-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_reports
ALTER COLUMN start_date DATE NULL;


ALTER TABLE risk_reports
ALTER COLUMN end_date DATE NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_reports
ALTER COLUMN start_date DATE NULL;


ALTER TABLE risk_reports
ALTER COLUMN end_date DATE NULL;
-- +goose StatementEnd
