-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_report_action_plans 
ADD proposed_controls NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_report_action_plans 
DROP COLUMN proposed_controls;
-- +goose StatementEnd
