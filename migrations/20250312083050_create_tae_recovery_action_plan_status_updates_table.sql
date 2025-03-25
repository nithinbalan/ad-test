-- +goose Up
-- +goose StatementBegin
CREATE TABLE tae_recovery_action_plan_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    test_and_exercise_recovery_action_plan_id INT NOT NULL,
    action_plan_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_tae_recovery_action_plan_status_updates_test_and_exercise_recovery_action_plans FOREIGN KEY (test_and_exercise_recovery_action_plan_id) REFERENCES test_and_exercise_recovery_action_plans (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_recovery_action_plan_status_updates_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_recovery_action_plan_status_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tae_recovery_action_plan_status_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tae_recovery_action_plan_status_updates;
-- +goose StatementEnd
