-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_report_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_report_id INT NOT NULL,
    compliance_report_requirement_id INT NOT NULL,
    title NVARCHAR(MAX) NOT NULL,
    description NVARCHAR(MAX) NOT NULL,
    action_plan_status_id INT NOT NULL,
    target_date DATETIME NOT NULL,
    severity_level_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_compliance_report_action_plans_compliance_reports FOREIGN KEY (compliance_report_id) REFERENCES compliance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_action_plans_compliance_report_requirements FOREIGN KEY (compliance_report_requirement_id) REFERENCES compliance_report_requirements (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_action_plans_users FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_action_plans_users_updated FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_report_action_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_report_action_plans;
-- +goose StatementEnd
