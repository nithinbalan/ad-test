-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_test_and_exercise_system_registers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    bcp_test_and_exercise_id INT NOT NULL,
    system_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_test_and_exercise_system_registers_test_exercise FOREIGN KEY (bcp_test_and_exercise_id) REFERENCES bcp_test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_test_and_exercise_system_registers_system_registers FOREIGN KEY (system_register_id) REFERENCES system_registers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_test_and_exercise_system_registers_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_test_and_exercise_system_registers;
-- +goose StatementEnd
