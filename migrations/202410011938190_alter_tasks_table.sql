-- +goose Up
-- +goose StatementBegin

-- Add the new risk_action_plan_id column
ALTER TABLE tasks
ADD audit_finding_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_audit_finding_action_plans FOREIGN KEY (audit_finding_action_plan_id) REFERENCES audit_finding_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the risk_action_plan_id column
ALTER TABLE taks
DROP CONSTRAINT FK_audit_finding_action_plans;

ALTER TABLE tasks
DROP COLUMN audit_finding_action_plan_id;
-- +goose StatementEnd
