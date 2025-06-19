-- +goose Up
-- +goose StatementBegin
ALTER TABLE grc_report_risk_controls
ALTER COLUMN title NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grc_report_risk_controls
ALTER COLUMN title NVARCHAR(255) NOT NULL;
-- +goose StatementEnd
