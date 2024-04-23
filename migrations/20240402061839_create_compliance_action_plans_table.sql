-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_action_plans (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    compliance_requirement_id INT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    severity_level_id INT NULL,
    target_date DATE NULL,
    responsible_user_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    comments NVARCHAR(MAX) NULL,
    action_plan_status_id INT NOT NULL,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_compliance_action_plans_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_users_2 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_action_plans_users_3 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_action_plans;
-- +goose StatementEnd
