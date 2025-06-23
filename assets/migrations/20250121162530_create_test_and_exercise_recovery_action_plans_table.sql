-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_recovery_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    test_and_exercise_id INT NOT NULL,
    bcp_phase_id INT NULL,
    process_id INT NOT NULL,
    title NVARCHAR(MAX) NULL,
    severity_level_id INT NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    target_completion NVARCHAR(MAX) NULL,
    target_completion_unit_id INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL, 
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_test_and_exercises FOREIGN KEY (test_and_exercise_id) REFERENCES test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_bcp_phases FOREIGN KEY (bcp_phase_id) REFERENCES bcp_phases(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_time_units FOREIGN KEY (target_completion_unit_id) REFERENCES time_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_recovery_action_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test_and_exercise_recovery_action_plans;
-- +goose StatementEnd
