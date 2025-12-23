-- +goose Up
-- +goose StatementBegin

UPDATE tenable_settings
SET status = CASE 
                WHEN is_active = 'true' THEN 'active'
                ELSE 'inactive'
            END;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE tenable_settings
SET status = NULL;

-- +goose StatementEnd
