-- +goose Up
-- +goose StatementBegin

-- Swap order: KPIs should be 3, Initiative should be 4
UPDATE strategy_profile_settings
SET [order] = CASE 
    WHEN title_en = 'Focus Area' THEN 1
    WHEN title_en = 'Objective' THEN 2
    WHEN title_en = 'KPIs' THEN 3
    WHEN title_en = 'Initiative' THEN 4
    WHEN title_en = 'Projects' THEN 5
    WHEN title_en = 'Milestone' THEN 6
    WHEN title_en = 'Action Plan' THEN 7
END
WHERE title_en IN (
    'Focus Area', 'Objective', 'Initiative', 'KPIs', 'Projects', 'Milestone', 'Action Plan'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Revert the order: KPIs back to 4, Initiative back to 3
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

-- +goose StatementEnd
