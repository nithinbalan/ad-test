-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_registers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    service_description NVARCHAR(MAX) NULL,
    modules_components NVARCHAR(MAX) NULL,
    third_party_dependencies NVARCHAR(MAX) NULL,
    service_importance_id INT NULL,
    responsible_it_team_id INT NULL,
    outsourcing_partner NVARCHAR(MAX) NULL,
    operation_location_id INT NULL,
    business_owner_id INT NULL,
    customers NVARCHAR(MAX) NULL,
    impact_type_id INT NULL,
    disruption_impact_level_id INT NULL,
    rto INT NULL,
    rto_unit_id INT NULL,
    mtpd INT NULL,
    mtpd_unit_id INT NULL,
    rpo INT NULL,
    rpo_unit_id INT NULL,
    it_processes NVARCHAR(MAX) NULL,
    it_service_owner_id INT NULL,
    process_operation_team_id INT NOT NULL,
    process_recovery_team_id INT NOT NULL,
    recovery_manager_id INT NULL,
    recovery_manager_phone VARCHAR(20) NOT NULL,
    recovery_manager_email VARCHAR(255) NOT NULL,
    recovery_location NVARCHAR(MAX) NULL,
    vendor_rto INT NULL,
    vendor_rto_unit_id INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_system_registers_service_importances FOREIGN KEY (service_importance_id) REFERENCES system_register_service_importances(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_job_teams_1 FOREIGN KEY (responsible_it_team_id) REFERENCES job_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_locations FOREIGN KEY (operation_location_id) REFERENCES locations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_organization_units FOREIGN KEY (business_owner_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_impact_categories FOREIGN KEY (impact_type_id) REFERENCES impact_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_disruption_levels FOREIGN KEY (disruption_impact_level_id) REFERENCES system_register_disruption_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_mtpd_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_designations_1 FOREIGN KEY (it_service_owner_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_job_teams_2 FOREIGN KEY (process_operation_team_id) REFERENCES job_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_job_teams_3 FOREIGN KEY (process_recovery_team_id) REFERENCES job_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_designations_2 FOREIGN KEY (recovery_manager_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_vendor_rto_time_units FOREIGN KEY (vendor_rto_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE system_registers;
-- +goose StatementEnd
