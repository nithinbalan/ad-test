-- +goose Up
-- +goose StatementBegin

-- Add focus_area_id column
ALTER TABLE strategy_profile_focus_areas
ADD focus_area_id INT NULL;

-- Remove old columns
ALTER TABLE strategy_profile_focus_areas
DROP COLUMN title,
            description;

-- Add foreign key constraint
ALTER TABLE strategy_profile_focus_areas
ADD CONSTRAINT FK_strategy_profile_focus_areas_focus_area FOREIGN KEY (focus_area_id)
    REFERENCES focus_areas(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Drop the foreign key
ALTER TABLE strategy_profile_focus_areas
DROP CONSTRAINT FK_strategy_profile_focus_areas_focus_area;

-- Add back the columns (SQL Server does NOT support AFTER)
ALTER TABLE strategy_profile_focus_areas
ADD title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL;

-- Drop the new column
ALTER TABLE strategy_profile_focus_areas
DROP COLUMN focus_area_id;

-- +goose StatementEnd
