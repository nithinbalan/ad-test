-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profile_projects (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategy_profile_id INT NOT NULL,
    strategy_profile_focus_area_id INT NULL,
    strategy_profile_goal_id INT NULL,
    strategy_profile_objective_id INT NULL,
    strategy_profile_initiative_id INT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_strategy_profile_projects_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
        REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
        REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_objectives FOREIGN KEY (strategy_profile_objective_id) 
        REFERENCES strategy_profile_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_initiatives FOREIGN KEY (strategy_profile_initiative_id) 
        REFERENCES strategy_profile_initiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_1 FOREIGN KEY (responsible_user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_2 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_3 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategy_profile_projects;
-- +goose StatementEnd
