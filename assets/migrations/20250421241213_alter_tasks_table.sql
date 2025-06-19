-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD strategy_profile_action_plan_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN strategy_profile_action_plan_id;
-- +goose StatementEnd
