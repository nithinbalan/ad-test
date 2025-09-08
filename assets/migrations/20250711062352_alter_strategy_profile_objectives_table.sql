-- +goose Up
-- +goose StatementBegin

-- Step 1: Rename the table
EXEC sp_rename 'strategy_profile_objectives', 'strategy_profile_strategic_objectives';

-- Step 2: Drop old foreign key constraints (since they reference the old table name)
ALTER TABLE strategy_profile_strategic_objectives
DROP CONSTRAINT FK_strategy_profile_objectives_strategy_profiles,
                 FK_strategy_profile_objectives_strategy_profile_focus_areas,
                 FK_strategy_profile_objectives_strategy_profile_goals,
                 FK_strategy_profile_objectives_users_1,
                 FK_strategy_profile_objectives_users_2;

-- Step 3: Add new columns
ALTER TABLE strategy_profile_strategic_objectives
ADD strategic_objective_id INT NULL,
    score DECIMAL(5,2) NULL;

-- Step 4: Drop old columns
ALTER TABLE strategy_profile_strategic_objectives
DROP COLUMN title,
            description;

-- Step 5: Recreate all foreign key constraints with updated names
ALTER TABLE strategy_profile_strategic_objectives
ADD CONSTRAINT FK_strategy_profile_strategic_objectives_strategy_profiles FOREIGN KEY (strategy_profile_id) 
    REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_objectives_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
    REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_objectives_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
    REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_objectives_users_1 FOREIGN KEY (created_by_id) 
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_objectives_users_2 FOREIGN KEY (updated_by_id) 
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_objectives_strategic_objectives FOREIGN KEY (strategic_objective_id) 
    REFERENCES strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop all foreign keys
ALTER TABLE strategy_profile_strategic_objectives
DROP CONSTRAINT FK_strategy_profile_strategic_objectives_strategy_profiles,
                 FK_strategy_profile_strategic_objectives_strategy_profile_focus_areas,
                 FK_strategy_profile_strategic_objectives_strategy_profile_goals,
                 FK_strategy_profile_strategic_objectives_users_1,
                 FK_strategy_profile_strategic_objectives_users_2,
                 FK_strategy_profile_strategic_objectives_strategic_objectives;

-- Step 2: Add back the dropped columns
ALTER TABLE strategy_profile_strategic_objectives
ADD title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL;

-- Step 3: Drop newly added columns
ALTER TABLE strategy_profile_strategic_objectives
DROP COLUMN strategic_objective_id,
             score;

-- Step 4: Rename the table back
EXEC sp_rename 'strategy_profile_strategic_objectives', 'strategy_profile_objectives';

-- Step 5: Re-create original foreign key constraints
ALTER TABLE strategy_profile_objectives
ADD CONSTRAINT FK_strategy_profile_objectives_strategy_profiles FOREIGN KEY (strategy_profile_id) 
    REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_objectives_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
    REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_objectives_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
    REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_objectives_users_1 FOREIGN KEY (created_by_id) 
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_objectives_users_2 FOREIGN KEY (updated_by_id) 
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
