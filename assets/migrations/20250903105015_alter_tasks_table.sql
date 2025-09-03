-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD threat_action_plan_id INT NULL;
 
ALTER TABLE tasks
ADD vulnerability_action_plan_id INT NULL;
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN threat_action_plan_id;
 
ALTER TABLE tasks
DROP COLUMN vulnerability_action_plan_id;
-- +goose StatementEnd