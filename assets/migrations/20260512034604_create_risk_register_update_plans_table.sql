-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_register_update_plans
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    division_id INT NOT NULL,
    business_unit_id INT NOT NULL,
    risk_register_update_plan_status_id INT NOT NULL,
    number_of_operating_entities INT NULL,
    workdays INT NULL,
    risk_champion_id INT NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    notes NVARCHAR(MAX) NULL,
    total_percentage_completed INT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    submitted_by_id INT NULL,
    approved_by_id INT NULL,
    submitted_at DATETIME NULL,
    approved_at DATETIME NULL,
    next_review_level INT NULL,
    workflow_status_id INT NULL,
    CONSTRAINT FK_risk_register_update_plans_divisions FOREIGN KEY (division_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_business_units FOREIGN KEY (business_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_risk_register_update_plan_statuses FOREIGN KEY (risk_register_update_plan_status_id) REFERENCES risk_register_update_plan_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_users_1 FOREIGN KEY (risk_champion_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_users_4 FOREIGN KEY (submitted_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_users_5 FOREIGN KEY (approved_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_register_update_plans_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_register_update_plans;
-- +goose StatementEnd