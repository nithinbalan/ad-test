-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD identity_review_action_plan_id INT NULL;

ALTER TABLE tasks
ADD CONSTRAINT FK_tasks_identity_review_action_plans
FOREIGN KEY (identity_review_action_plan_id)
REFERENCES identity_review_action_plans(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP CONSTRAINT FK_tasks_identity_review_action_plans;

ALTER TABLE tasks
DROP COLUMN identity_review_action_plan_id;
-- +goose StatementEnd
