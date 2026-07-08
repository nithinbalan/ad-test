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
