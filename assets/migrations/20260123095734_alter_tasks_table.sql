-- +goose Up
-- +goose StatementBegin
-- Add the new violation_report_action_plan_id column
ALTER TABLE tasks
ADD violation_report_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_violation_report_action_plans FOREIGN KEY (violation_report_action_plan_id) REFERENCES violation_report_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the violation_report_action_plan_id column
ALTER TABLE tasks
DROP CONSTRAINT FK_violation_report_action_plans;

ALTER TABLE tasks
DROP COLUMN violation_report_action_plan_id;
-- +goose StatementEnd
