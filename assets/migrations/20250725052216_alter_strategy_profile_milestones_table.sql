-- +goose Up
-- +goose StatementBegin

-- Step 1: Rename the table
EXEC sp_rename 'strategy_profile_milestones', 'strategy_profile_project_milestones';

-- Step 2: Drop old foreign key constraints
ALTER TABLE strategy_profile_project_milestones
DROP CONSTRAINT FK_strategy_profile_milestones_strategy_profiles,
                 FK_strategy_profile_milestones_strategy_profile_focus_areas,
                 FK_strategy_profile_milestones_strategy_profile_goals,
                 FK_strategy_profile_milestones_strategy_profile_objectives,
                 FK_strategy_profile_milestones_strategy_profile_initiatives,
                 FK_strategy_profile_milestones_strategy_profile_projects,
                 FK_strategy_profile_milestones_users_1,
                 FK_strategy_profile_milestones_users_2;

-- Step 3: Add new columns
ALTER TABLE strategy_profile_project_milestones
ADD project_milestone_id INT NULL;

-- Step 4: Drop old columns
ALTER TABLE strategy_profile_project_milestones
DROP COLUMN title,
            strategy_profile_focus_area_id,
            strategy_profile_goal_id,
            strategy_profile_objective_id,
            strategy_profile_initiative_id;

-- Step 5: Recreate foreign keys with updated references
ALTER TABLE strategy_profile_project_milestones
ADD CONSTRAINT FK_strategy_profile_project_milestones_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_milestones_strategy_profile_projects FOREIGN KEY (strategy_profile_project_id) 
        REFERENCES strategy_profile_projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_milestones_project_milestones FOREIGN KEY (project_milestone_id) 
        REFERENCES project_milestones(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_milestones_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_project_milestones_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop all new foreign key constraints
ALTER TABLE strategy_profile_project_milestones
DROP CONSTRAINT FK_strategy_profile_project_milestones_strategy_profiles,
                 FK_strategy_profile_project_milestones_strategy_profile_projects,
                 FK_strategy_profile_project_milestones_project_milestones,
                 FK_strategy_profile_project_milestones_users_1,
                 FK_strategy_profile_project_milestones_users_2;


-- Step 2: Re-add dropped columns
ALTER TABLE strategy_profile_project_milestones
ADD title NVARCHAR(500) NULL,
    strategy_profile_objective_id INT NULL,
    strategy_profile_goal_id INT NULL,
    strategy_profile_focus_area_id INT NULL,
    strategy_profile_initiative_id INT NULL;

-- Step 3: Drop newly added columns
ALTER TABLE strategy_profile_project_milestones
DROP COLUMN project_milestone_id;

-- Step 4: Rename the table back
EXEC sp_rename 'strategy_profile_project_milestones', 'strategy_profile_milestones';

-- Step 5: Re-create original foreign key constraints
ALTER TABLE strategy_profile_milestones
ADD CONSTRAINT FK_strategy_profile_milestones_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
        REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
        REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_strategy_profile_objectives FOREIGN KEY (strategy_profile_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_strategy_profile_initiatives FOREIGN KEY (strategy_profile_initiative_id) 
        REFERENCES strategy_profile_strategic_intiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_strategy_profile_projects FOREIGN KEY (strategy_profile_project_id) 
        REFERENCES strategy_profile_projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_milestones_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd