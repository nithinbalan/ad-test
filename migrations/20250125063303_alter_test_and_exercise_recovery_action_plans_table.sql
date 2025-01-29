-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
DROP CONSTRAINT FK_test_and_exercise_recovery_action_plans_time_units ;

ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN target_completion_unit_id;

ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN target_completion;

ALTER TABLE test_and_exercise_recovery_action_plans ADD target_date DATE NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_recovery_action_plans
ADD target_completion NVARCHAR(MAX) NULL;

ALTER TABLE test_and_exercise_recovery_action_plans
ADD target_completion_unit_id INT NULL;

ALTER TABLE test_and_exercise_recovery_action_plans
ADD CONSTRAINT FK_test_and_exercise_recovery_action_plans_time_units FOREIGN KEY (target_completion_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE test_and_exercise_recovery_action_plans
DROP COLUMN target_date;
-- +goose StatementEnd
