-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_test_results (
    id INT PRIMARY KEY IDENTITY(1, 1),
    test_and_exercise_id INT NOT NULL,
    test_completion_date DATE NOT NULL,
    communication_feedback_option_id INT NOT NULL,
    recovery_feedback_option_id INT NOT NULL,
    test_result_option_id INT NOT NULL,
    test_observations NVARCHAR(MAX) NULL,
    recommendations NVARCHAR(MAX) NULL,
    findings NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_test_and_exercise_test_results_test_and_exercises FOREIGN KEY (test_and_exercise_id) REFERENCES test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_results_feedback_options_1 FOREIGN KEY (communication_feedback_option_id) REFERENCES feedback_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_results_feedback_options_2 FOREIGN KEY (recovery_feedback_option_id) REFERENCES feedback_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_results_test_result_options FOREIGN KEY (test_result_option_id) REFERENCES test_result_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_results_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_test_results_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.test_and_exercise_test_results;
-- +goose StatementEnd
