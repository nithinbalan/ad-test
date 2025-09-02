-- +goose Up
-- +goose StatementBegin

-- Step 1: Drop old foreign key constraints
ALTER TABLE strategy_profile_projects
DROP CONSTRAINT FK_strategy_profile_projects_strategy_profile_focus_areas,
                 FK_strategy_profile_projects_strategy_profile_goals,
                 FK_strategy_profile_projects_strategy_profile_objectives,
                 FK_strategy_profile_projects_strategy_profile_initiatives,
                 FK_strategy_profile_projects_users_1,
                 FK_strategy_profile_projects_users_2,
                 FK_strategy_profile_projects_users_3;

-- Step 2: Add new column
ALTER TABLE strategy_profile_projects
ADD project_id INT NULL,
    strategy_profile_strategic_objective_id INT NULL,
    strategy_profile_strategic_intiative_id INT NULL;


-- Step 3: Drop old columns
ALTER TABLE strategy_profile_projects
DROP COLUMN title,
            description,
            responsible_user_id,
            strategy_profile_objective_id,
            strategy_profile_goal_id,
            strategy_profile_focus_area_id,
            strategy_profile_initiative_id,
            start_date,
            end_date;

-- Step 4: Recreate foreign keys with updated references
ALTER TABLE strategy_profile_projects
ADD CONSTRAINT FK_strategy_profile_projects_projects FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_strategic_objectives FOREIGN KEY (strategy_profile_strategic_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_strategic_intiatives FOREIGN KEY (strategy_profile_strategic_intiative_id) 
        REFERENCES strategy_profile_strategic_intiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop new foreign key constraints
ALTER TABLE strategy_profile_projects
DROP CONSTRAINT FK_strategy_profile_projects_projects,
                 FK_strategy_profile_projects_strategy_profile_strategic_objectives,
                 FK_strategy_profile_projects_strategy_profile_strategic_intiatives,
                 FK_strategy_profile_projects_users_1,
                 FK_strategy_profile_projects_users_2;

-- Step 1.1: Drop users_3 if exists (safe rollback)
IF EXISTS (
    SELECT 1 FROM sys.foreign_keys 
    WHERE name = 'FK_strategy_profile_projects_users_3'
)
BEGIN
    ALTER TABLE strategy_profile_projects
    DROP CONSTRAINT FK_strategy_profile_projects_users_3;
END

-- Step 2: Re-add previously dropped columns
ALTER TABLE strategy_profile_projects
ADD title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    strategy_profile_objective_id INT NULL,
    strategy_profile_goal_id INT NULL,
    strategy_profile_focus_area_id INT NULL,
    strategy_profile_initiative_id INT NULL,
    start_date DATE NULL,
    end_date DATE NULL;

-- Step 3: Drop the newly added column
ALTER TABLE strategy_profile_projects
DROP COLUMN project_id,
            strategy_profile_strategic_objective_id,
            strategy_profile_strategic_intiative_id;

-- Step 4: Recreate original foreign key constraints
ALTER TABLE strategy_profile_projects
ADD CONSTRAINT FK_strategy_profile_projects_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
        REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
        REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_objectives FOREIGN KEY (strategy_profile_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_strategy_profile_initiatives FOREIGN KEY (strategy_profile_initiative_id) 
        REFERENCES strategy_profile_strategic_intiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_1 FOREIGN KEY (responsible_user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_2 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_projects_users_3 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
