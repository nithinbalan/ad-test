-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD stakeholder_action_plan_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN stakeholder_action_plan_id;
-- +goose StatementEnd
