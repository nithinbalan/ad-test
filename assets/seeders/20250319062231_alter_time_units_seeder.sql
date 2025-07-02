-- +goose Up
-- +goose StatementBegin
UPDATE time_units 
SET type = CASE 
    WHEN title_en = 'Hours' THEN 'hour'
    WHEN title_en = 'Minutes' THEN 'minute'
    WHEN title_en = 'Days' THEN 'day'
    WHEN title_en = 'Weeks' THEN 'week'
    ELSE NULL 
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE time_units 
SET type = NULL;
-- +goose StatementEnd
