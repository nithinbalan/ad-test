-- +goose Up
-- +goose StatementBegin
UPDATE action_plan_statuses
SET 
    title_en = 'Closed',
    title_ar = 'مغلق'
WHERE 
    type = 'close';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE action_plan_statuses
SET 
    title_en = 'Close',
    title_ar = 'مغلق'
WHERE 
    type = 'close';
-- +goose StatementEnd
