-- +goose Up
-- +goose StatementBegin
-- Add the new project_action_plan_id column
ALTER TABLE tasks
ADD project_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_project_action_plans FOREIGN KEY (project_action_plan_id) REFERENCES project_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the project_action_plan_id column
ALTER TABLE tasks
DROP CONSTRAINT FK_project_action_plans;

ALTER TABLE tasks
DROP COLUMN project_action_plan_id;
-- +goose StatementEnd