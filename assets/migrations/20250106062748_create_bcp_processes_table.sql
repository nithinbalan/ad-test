-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_continuity_plan_id INT NULL,
    process_id INT NULL,
    risk_level_id INT NULL,
    rto INt NULL,
    rto_unit_id INt NULL,
    mtpd INt NULL,
    mtpd_unit_id INt NULL,
    rpo INt NULL,
    rpo_unit_id INt NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,
    CONSTRAINT FK_bcp_processes_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id),
    CONSTRAINT FK_bcp_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id),
    CONSTRAINT FK_bcp_processes_risk_level FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_processes_time_units_1 FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_processes_time_units_2 FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_processes_time_units_3 FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
    CONSTRAINT FK_bcp_processes_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_processes;
-- +goose StatementEnd
