-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
ADD communication_action_plan_status_id INT NULL;

ALTER TABLE test_and_exercise_communication_action_plans ADD description NVARCHAR(MAX) NULL;

ALTER TABLE test_and_exercise_communication_action_plans
ADD CONSTRAINT FK_test_and_exercise_communication_action_plans_communication_action_plan_statuses FOREIGN KEY (communication_action_plan_status_id) REFERENCES communication_action_plan_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
DROP CONSTRAINT FK_test_and_exercise_communication_action_plans_communication_action_plan_statuses;

ALTER TABLE test_and_exercise_communication_action_plans
DROP COLUMN communication_action_plan_status_id;

ALTER TABLE test_and_exercise_communication_action_plans 
DROP COLUMN description;
-- +goose StatementEnd
