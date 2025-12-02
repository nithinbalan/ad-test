-- +goose Up
-- +goose NO TRANSACTION

UPDATE tenable_settings
SET status = CASE 
                WHEN is_active = 'true' THEN 'active'
                ELSE 'inactive'
            END;
