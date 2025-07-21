-- +goose Up
-- +goose StatementBegin

-- Drop the existing foreign key constraint
ALTER TABLE project_action_plans
DROP CONSTRAINT FK_project_action_plans_action_plan_statuses;

-- Alter the column to allow NULL
ALTER TABLE project_action_plans
ALTER COLUMN action_plan_status_id INT NULL;

-- Re-add the foreign key constraint (with NULL support)
ALTER TABLE project_action_plans
ADD CONSTRAINT FK_project_action_plans_action_plan_statuses 
    FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the nullable foreign key constraint
ALTER TABLE project_action_plans
DROP CONSTRAINT FK_project_action_plans_action_plan_statuses;

-- Revert the column to NOT NULL
ALTER TABLE project_action_plans
ALTER COLUMN action_plan_status_id INT NOT NULL;

-- Re-add the original NOT NULL foreign key constraint
ALTER TABLE project_action_plans
ADD CONSTRAINT FK_project_action_plans_action_plan_statuses 
    FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
