-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN description NVARCHAR(MAX) NULL;
ALTER TABLE compliance_report_action_plans
ALTER COLUMN target_date DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN description NVARCHAR(MAX) NOT NULL;
ALTER TABLE compliance_report_action_plans
ALTER COLUMN target_date DATETIME NOT NULL;
-- +goose StatementEnd