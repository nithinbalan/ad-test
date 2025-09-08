-- +goose Up
-- +goose StatementBegin

-- Step 1: Rename the table
EXEC sp_rename 'strategy_profile_initiatives', 'strategy_profile_strategic_intiatives';

-- Step 2: Drop old foreign key constraints
ALTER TABLE strategy_profile_strategic_intiatives
DROP CONSTRAINT FK_strategy_profile_initiatives_strategy_profiles,
                 FK_strategy_profile_initiatives_strategy_profile_focus_areas,
                 FK_strategy_profile_initiatives_strategy_profile_goals,
                 FK_strategy_profile_initiatives_strategy_profile_objectives,
                 FK_strategy_profile_initiatives_users_1,
                 FK_strategy_profile_initiatives_users_2,
                 FK_strategy_profile_initiatives_users_3;

-- Step 3: Add new columns
ALTER TABLE strategy_profile_strategic_intiatives
ADD strategic_intiative_id INT NULL,
    strategy_profile_strategic_objective_id INT NULL;

-- Step 4: Drop old columns
ALTER TABLE strategy_profile_strategic_intiatives
DROP COLUMN title,
            description,
            responsible_user_id,
            strategy_profile_objective_id,
            strategy_profile_goal_id,
            strategy_profile_focus_area_id,
            start_date,
            end_date;

-- Step 5: Recreate foreign keys with updated references
ALTER TABLE strategy_profile_strategic_intiatives
ADD CONSTRAINT FK_strategy_profile_strategic_intiatives_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_intiatives_strategy_profile_strategic_objectives FOREIGN KEY (strategy_profile_strategic_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_intiatives_strategic_intiatives FOREIGN KEY (strategic_intiative_id) 
        REFERENCES strategic_intiatives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_intiatives_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_strategic_intiatives_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop all new foreign key constraints
ALTER TABLE strategy_profile_strategic_intiatives
DROP CONSTRAINT FK_strategy_profile_strategic_intiatives_strategy_profiles,
                 FK_strategy_profile_strategic_intiatives_strategy_profile_strategic_objectives,
                 FK_strategy_profile_strategic_intiatives_strategic_intiatives,
                 FK_strategy_profile_strategic_intiatives_users_1,
                 FK_strategy_profile_strategic_intiatives_users_2;

-- Step 1.1: Drop users_3 (missed previously)
IF EXISTS (
    SELECT 1 FROM sys.foreign_keys 
    WHERE name = 'FK_strategy_profile_strategic_intiatives_users_3'
)
BEGIN
    ALTER TABLE strategy_profile_strategic_intiatives
    DROP CONSTRAINT FK_strategy_profile_strategic_intiatives_users_3;
END

-- Step 2: Re-add dropped columns
ALTER TABLE strategy_profile_strategic_intiatives
ADD title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    strategy_profile_objective_id INT NULL,
    strategy_profile_goal_id INT NULL,
    strategy_profile_focus_area_id INT NULL,
    start_date DATE NULL,
    end_date DATE NULL;

-- Step 3: Drop newly added columns
ALTER TABLE strategy_profile_strategic_intiatives
DROP COLUMN strategic_intiative_id,
            strategy_profile_strategic_objective_id;

-- Step 4: Rename the table back
EXEC sp_rename 'strategy_profile_strategic_intiatives', 'strategy_profile_initiatives';

-- Step 5: Re-create original foreign key constraints
ALTER TABLE strategy_profile_initiatives
ADD CONSTRAINT FK_strategy_profile_initiatives_strategy_profiles FOREIGN KEY (strategy_profile_id) 
        REFERENCES strategy_profiles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_strategy_profile_focus_areas FOREIGN KEY (strategy_profile_focus_area_id) 
        REFERENCES strategy_profile_focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_strategy_profile_goals FOREIGN KEY (strategy_profile_goal_id) 
        REFERENCES strategy_profile_goals(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_strategy_profile_objectives FOREIGN KEY (strategy_profile_objective_id) 
        REFERENCES strategy_profile_strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_users_1 FOREIGN KEY (responsible_user_id) 
        REFERENCES users(id) ON DELETE SET NULL ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_users_2 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_initiatives_users_3 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd