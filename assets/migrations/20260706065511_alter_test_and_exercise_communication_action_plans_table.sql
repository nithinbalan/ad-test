-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
ALTER COLUMN process_id INT NULL;

ALTER TABLE test_and_exercise_communication_action_plans
ADD system_register_id INT NULL;

ALTER TABLE test_and_exercise_communication_action_plans
ADD CONSTRAINT FK_test_and_exercise_communication_action_plans_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_communication_action_plans
DROP CONSTRAINT FK_test_and_exercise_communication_action_plans_system_registers;

ALTER TABLE test_and_exercise_communication_action_plans
DROP COLUMN system_register_id;

ALTER TABLE test_and_exercise_communication_action_plans
ALTER COLUMN process_id INT NOT NULL;
-- +goose StatementEnd
