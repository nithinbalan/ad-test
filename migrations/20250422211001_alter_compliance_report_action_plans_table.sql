-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN severity_level_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_report_action_plans
ALTER COLUMN severity_level_id INT NOT NULL;
-- +goose StatementEnd