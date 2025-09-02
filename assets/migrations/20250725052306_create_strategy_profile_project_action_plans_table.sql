-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profile_project_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategy_profile_id INT NOT NULL,
    strategy_profile_project_id INT NULL,
    strategy_profile_project_milestone_id INT NULL,
    project_action_plan_id INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_strategy_profile_project_action_plans_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_action_plans_strategy_profile_projects FOREIGN KEY (strategy_profile_project_id) 
        REFERENCES strategy_profile_projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_action_plans_strategy_profile_project_milestones FOREIGN KEY (strategy_profile_project_milestone_id) 
        REFERENCES strategy_profile_project_milestones(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_action_plans_project_action_plans FOREIGN KEY (project_action_plan_id) 
        REFERENCES project_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_action_plans_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_action_plans_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  dbo.strategy_profile_project_action_plans;
-- +goose StatementEnd
