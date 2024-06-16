-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    module_id INT NOT NULL,
    compliance_action_plan_id INT NULL,
    sla_action_plan_id INT NULL,
    incident_action_plan_id INT NULL,
    control_action_plan_id INT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    target_date DATE NULL,
    url VARCHAR(1000) NOT NULL,
    severity_level_id INT NULL,
    responsible_user_id INT NULL,
    task_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_tasks_modules FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_compliance_action_plans FOREIGN KEY (compliance_action_plan_id) REFERENCES compliance_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_sla_action_plans FOREIGN KEY (sla_action_plan_id) REFERENCES sla_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_incident_action_plans FOREIGN KEY (incident_action_plan_id) REFERENCES incident_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_control_action_plans FOREIGN KEY (control_action_plan_id) REFERENCES control_action_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_task_statuses FOREIGN KEY (task_status_id) REFERENCES task_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_tasks_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;

-- +goose StatementEnd