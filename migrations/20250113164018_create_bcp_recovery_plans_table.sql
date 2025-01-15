-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_recovery_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_continuity_plan_id INT NOT NULL,
    bcp_phase_id INT NOT NULL,
    bcp_process_id INT NOT NULL,
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
   CONSTRAINT FK_bcp_recovery_plans_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_bcp_phases FOREIGN KEY (bcp_phase_id) REFERENCES bcp_phases(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_bcp_processes FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_organization_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_time_units FOREIGN KEY (target_completion_unit_id) REFERENCES time_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_recovery_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.bcp_recovery_plans;
-- +goose StatementEnd
