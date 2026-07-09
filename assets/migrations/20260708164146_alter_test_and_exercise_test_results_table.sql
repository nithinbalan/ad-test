-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results DROP CONSTRAINT FK_test_and_exercise_test_results_feedback_options_1;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results DROP CONSTRAINT FK_test_and_exercise_test_results_feedback_options_2;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results DROP CONSTRAINT FK_test_and_exercise_test_results_test_result_options;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results DROP COLUMN test_completion_date, communication_feedback_option_id, recovery_feedback_option_id, test_result_option_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results ADD
    test_completion_date DATE  NULL,
    communication_feedback_option_id INT NULL,
    recovery_feedback_option_id INT NULL,
    test_result_option_id INT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results
ADD CONSTRAINT FK_test_and_exercise_test_results_feedback_options_1
FOREIGN KEY (communication_feedback_option_id)
REFERENCES feedback_options(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results
ADD CONSTRAINT FK_test_and_exercise_test_results_feedback_options_2
FOREIGN KEY (recovery_feedback_option_id)
REFERENCES feedback_options(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE test_and_exercise_test_results
ADD CONSTRAINT FK_test_and_exercise_test_results_test_result_options
FOREIGN KEY (test_result_option_id)
REFERENCES test_result_options(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd