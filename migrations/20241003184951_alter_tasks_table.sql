-- +goose Up
-- +goose StatementBegin
-- Add the new meeting_action_plan_id column
ALTER TABLE tasks
ADD meeting_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_meeting_action_plans FOREIGN KEY (meeting_action_plan_id) REFERENCES meeting_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the meeting_action_plan_id column
ALTER TABLE taks
DROP CONSTRAINT FK_meeting_action_plans;

ALTER TABLE tasks
DROP COLUMN meeting_finding_action_plan_id;
-- +goose StatementEnd
