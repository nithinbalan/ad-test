-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
ADD action_plan_status_id INT NULL;

ALTER TABLE test_and_exercise_recovery_action_plans
ADD CONSTRAINT FK_test_and_exercise_recovery_action_plans_action_plan_statuses
FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
DROP CONSTRAINT FK_test_and_exercise_recovery_action_plans_action_plan_statuses;

ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN action_plan_status_id;
-- +goose StatementEnd
