-- +goose Up
-- +goose StatementBegin

-- Step 1: Assign order to the correct items (excluding Goals and Target)
UPDATE strategy_profile_settings 
SET [order] = CASE 
    WHEN title_en = 'Focus Area' THEN 1
    WHEN title_en = 'Objective' THEN 2
    WHEN title_en = 'Initiative' THEN 3
    WHEN title_en = 'KPIs' THEN 4
    WHEN title_en = 'Projects' THEN 5
    WHEN title_en = 'Milestone' THEN 6
    WHEN title_en = 'Action Plan' THEN 7
END
WHERE title_en IN (
    'Focus Area', 'Objective', 'Initiative', 'KPIs', 'Projects', 'Milestone', 'Action Plan'
);

-- Step 2: Remove "Goals" and "Target"
DELETE FROM strategy_profile_settings 
WHERE title_en IN ('Goals', 'Target');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Step 1: Re-insert "Goals" and "Target"
INSERT INTO strategy_profile_settings (title_en, title_ar, is_enable, type)
VALUES 
('Goals', 'الأهداف', 1, 'goal'),
('Target', 'الهدف', 1, 'target');

-- Step 2: Reset order to NULL for all
UPDATE strategy_profile_settings 
SET [order] = NULL;

-- +goose StatementEnd
