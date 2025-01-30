-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    test_and_exercise_id INT NOT NULL,
    process_id INT NOT NULL,
    planned_rto INT NULL,
    planned_rto_unit_id INT NULL,
    planned_mtpd INT NULL,
    planned_mtpd_unit_id INT NULL,
    planned_rpo INT NULL,
    planned_rpo_unit_id INT NULL,
    actual_rto INT NULL,
    actual_rto_unit_id INT NULL,
    actual_mtpd INT NULL,
    actual_mtpd_unit_id INT NULL,
    actual_rpo INT NULL,
    actual_rpo_unit_id INT NULL,
    comments VARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_test_and_exercise_processes_test_and_exercises FOREIGN KEY (test_and_exercise_id) REFERENCES test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_1 FOREIGN KEY (planned_rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_2 FOREIGN KEY (planned_mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_3 FOREIGN KEY (planned_rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_4 FOREIGN KEY (actual_rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_5 FOREIGN KEY (actual_mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_time_units_6 FOREIGN KEY (actual_rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_processes_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.test_and_exercise_processes;
-- +goose StatementEnd
