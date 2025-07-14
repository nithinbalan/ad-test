-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_action_plans
ALTER COLUMN project_milestone_id INT NULL;

ALTER TABLE project_action_plans
ADD project_id INT NULL;

ALTER TABLE project_action_plans
ADD CONSTRAINT FK_project_action_plans_projects FOREIGN KEY (project_id)
REFERENCES projects (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_action_plans
DROP CONSTRAINT FK_project_action_plans_projects;

ALTER TABLE project_action_plans
DROP COLUMN project_id;

ALTER TABLE project_action_plans
ALTER COLUMN project_milestone_id INT NOT NULL;
-- +goose StatementEnd


