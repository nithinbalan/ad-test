-- +goose Up
-- +goose StatementBegin

-- Add the new boolean field 'is_closed'
ALTER TABLE action_plan_statuses
ADD is_closed BIT NOT NULL DEFAULT 0;

-- Drop the old constraint on the 'type' field (if needed)
ALTER TABLE action_plan_statuses
DROP CONSTRAINT CHK_action_plan_statuses_type;

-- Modify the 'type' field to accept both the old and new values
ALTER TABLE action_plan_statuses
ADD CONSTRAINT CHK_action_plan_statuses_type
CHECK (type IN ('new', 'in_progress', 'completed', 'on_hold', 'terminated', 'close', 'reopen'));

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Remove the 'is_closed' field
ALTER TABLE action_plan_statuses
DROP COLUMN is_closed;

-- Drop the modified constraint on the 'type' field
ALTER TABLE action_plan_statuses
DROP CONSTRAINT CHK_action_plan_statuses_type;

-- Restore the original 'type' constraint with old values
ALTER TABLE action_plan_statuses
ADD CONSTRAINT CHK_action_plan_statuses_type
CHECK (type IN ('new', 'in_progress', 'completed', 'on_hold', 'terminated'));

-- +goose StatementEnd
