-- +goose Up
-- +goose StatementBegin
UPDATE action_plan_statuses 
SET color_code = '#FD7E14' 
WHERE type = 'new';

UPDATE action_plan_statuses 
SET color_code = '#FFC107' 
WHERE type = 'in_progress';

UPDATE action_plan_statuses 
SET color_code = '#90DD75' 
WHERE type = 'completed';

UPDATE action_plan_statuses 
SET color_code = '#28A745' 
WHERE type = 'close';

UPDATE action_plan_statuses 
SET color_code = '#6C757D' 
WHERE type = 'on_hold';

UPDATE action_plan_statuses 
SET color_code = '#343A40' 
WHERE type = 'terminated';

UPDATE action_plan_statuses 
SET color_code = '#DC3545' 
WHERE type = 'reopen';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE action_plan_statuses 
SET color_code = '#007BFF' 
WHERE type = 'new';

UPDATE action_plan_statuses 
SET color_code = '#17A2B8' 
WHERE type = 'in_progress';

UPDATE action_plan_statuses 
SET color_code = '#28A745' 
WHERE type = 'completed';

UPDATE action_plan_statuses 
SET color_code = '#6C757D' 
WHERE type = 'close';

UPDATE action_plan_statuses 
SET color_code = '#FFC107' 
WHERE type = 'on_hold';

UPDATE action_plan_statuses 
SET color_code = '#DC3545' 
WHERE type = 'terminated';

UPDATE action_plan_statuses 
SET color_code = '#6610F2' 
WHERE type = 'reopen';
-- +goose StatementEnd
