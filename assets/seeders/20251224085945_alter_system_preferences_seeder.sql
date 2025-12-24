-- +goose Up
-- +goose StatementBegin

UPDATE system_preferences
SET theme_id = (
    SELECT TOP 1 id
    FROM themes
    WHERE is_system = 1
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE system_preferences
SET theme_id = NULL;

-- +goose StatementEnd
