-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_action_plans
ADD severity_level_id INT NULL,
    review_user_id INT NULL;

ALTER TABLE project_action_plans
ADD CONSTRAINT FK_project_action_plans_severity_levels 
    FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE project_action_plans
ADD CONSTRAINT FK_project_action_plans_users_4 
    FOREIGN KEY (review_user_id) REFERENCES users(id) 
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_action_plans
DROP CONSTRAINT FK_project_action_plans_severity_levels;

ALTER TABLE project_action_plans
DROP CONSTRAINT FK_project_action_plans_users_4;

ALTER TABLE project_action_plans
DROP COLUMN severity_level_id,
             review_user_id;
-- +goose StatementEnd

