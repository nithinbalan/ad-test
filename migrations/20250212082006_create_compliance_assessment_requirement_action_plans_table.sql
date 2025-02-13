-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirement_action_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_requirement_id INT NOT NULL,
    title NVARCHAR(MAX) NULL,
    responsible_user_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    target_date DATE NULL,
    description NVARCHAR(MAX) NULL,
    percentage INT NOT NULL DEFAULT 0,
    action_plan_status_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL, 
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_compliance_assessment_requirements FOREIGN KEY (compliance_assessment_requirement_id) REFERENCES compliance_assessment_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_action_plan_statuses FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_action_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_requirement_action_plans;
-- +goose StatementEnd
