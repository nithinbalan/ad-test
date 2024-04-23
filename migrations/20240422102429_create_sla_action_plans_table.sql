-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_action_plans (
    id INT PRIMARY KEY IDENTITY(1, 1),
    sla_register_id INT NOT NULL,
    sla_check_question_id INT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    severity_level_id INT NULL,
    target_date DATE NULL,
    responsible_user_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    comments NVARCHAR(MAX) NULL,
    action_plan_status_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_sla_action_plans_sla_register FOREIGN KEY (sla_register_id) REFERENCES sla_register (id),
    CONSTRAINT FK_sla_action_plans_sla_check_questions FOREIGN KEY (sla_check_question_id) REFERENCES sla_check_questions (id),
    CONSTRAINT FK_sla_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id),
    CONSTRAINT FK_sla_action_plans_users1 FOREIGN KEY (responsible_user_id) REFERENCES users (id),
    CONSTRAINT FK_sla_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id),
    CONSTRAINT FK_sla_action_plans_users2 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_sla_action_plans_users3 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sla_action_plans;
-- +goose StatementEnd
