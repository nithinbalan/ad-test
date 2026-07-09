-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_test_result_system_registers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    test_and_exercise_test_result_id INT NOT NULL,
    system_register_id INT NOT NULL,
    test_completion_date DATE NULL,
    communication_feedback_option_id INT NULL,
    recovery_feedback_option_id INT NULL,
    test_result_option_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_tae_test_result_system_registers_test_results FOREIGN KEY (test_and_exercise_test_result_id) REFERENCES test_and_exercise_test_results(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_test_result_system_registers_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_test_result_system_registers_feedback_options_1 FOREIGN KEY (communication_feedback_option_id) REFERENCES feedback_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_test_result_system_registers_feedback_options_2 FOREIGN KEY (recovery_feedback_option_id) REFERENCES feedback_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_test_result_system_registers_test_result_options FOREIGN KEY (test_result_option_id) REFERENCES test_result_options(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_test_result_system_registers_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test_and_exercise_test_result_system_registers;
-- +goose StatementEnd
