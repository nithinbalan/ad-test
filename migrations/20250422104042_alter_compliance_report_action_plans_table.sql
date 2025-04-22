-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN description NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN description NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd