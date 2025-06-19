-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_test_and_exercise_bcp_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    bcp_test_and_exercise_id INT NOT NULL,
    bcp_process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_test_and_exercise_bcp_processes_test_exercise FOREIGN KEY (bcp_test_and_exercise_id) REFERENCES bcp_test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_test_and_exercise_bcp_processes FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_test_and_exercise_bcp_processes;
-- +goose StatementEnd
