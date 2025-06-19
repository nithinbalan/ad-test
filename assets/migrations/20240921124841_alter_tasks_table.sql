-- +goose Up
-- +goose StatementBegin

-- Add the new risk_action_plan_id column
ALTER TABLE tasks
ADD risk_action_plan_id INT NULL;

-- Remove the existing url column
ALTER TABLE tasks
DROP COLUMN url;

-- Add the new url_param and query_param columns
ALTER TABLE tasks
ADD url_param NVARCHAR(MAX) NULL,
    query_param NVARCHAR(MAX) NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the risk_action_plan_id column
ALTER TABLE tasks
DROP COLUMN risk_action_plan_id;

-- Add the url column back
ALTER TABLE tasks
ADD url VARCHAR(1000) NOT NULL;

-- Remove the url_param and query_param columns
ALTER TABLE tasks
DROP COLUMN url_param,
    DROP COLUMN query_param;

-- +goose StatementEnd
