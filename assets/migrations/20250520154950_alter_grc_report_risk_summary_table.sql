-- +goose Up
-- +goose StatementBegin
ALTER TABLE grc_report_risk_summary
ADD risk_trend NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grc_report_risk_summary
DROP COLUMN risk_trend;
-- +goose StatementEnd
