-- +goose Up
-- +goose StatementBegin
UPDATE strategy_profile_settings 
SET type = CASE 
    WHEN title_en = 'Focus Area' THEN 'focus-area'
    WHEN title_en = 'Goals' THEN 'goal'
    WHEN title_en = 'Objective' THEN 'objective'
    WHEN title_en = 'Initiative' THEN 'initiative'
    WHEN title_en = 'Projects' THEN 'project'
    WHEN title_en = 'Milestone' THEN 'milestone'
    WHEN title_en = 'Action Plan' THEN 'action-plan'
    WHEN title_en = 'KPIs' THEN 'kpi'
    WHEN title_en = 'Target' THEN 'target'
END
WHERE title_en IN (
    'Focus Area', 'Goals', 'Objective', 'Initiative', 'Projects', 'Milestone', 'Action Plan', 'KPIs', 'Target'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE strategy_profile_settings 
SET type = NULL;
-- +goose StatementEnd
