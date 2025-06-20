-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_and_exercise_communication_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    test_and_exercise_id INT NOT NULL,
    bcp_phase_id INT NULL,
    process_id INT NOT NULL,
    title NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    responsible_unit_id INT NULL,
    recipient NVARCHAR(MAX) NULL,
    frequency NVARCHAR(MAX) NULL,
    channel NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_test_and_exercises FOREIGN KEY (test_and_exercise_id) REFERENCES test_and_exercises(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_bcp_phases FOREIGN KEY (bcp_phase_id) REFERENCES bcp_phases(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_test_and_exercise_communication_action_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test_and_exercise_communication_action_plans;
-- +goose StatementEnd
