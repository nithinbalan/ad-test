-- +goose Up
-- +goose StatementBegin

CREATE TABLE system_registers (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(100) NOT NULL UNIQUE,  
    number_of_people_necessary INT NULL,
    business_owner_unit_id INT NULL,
    it_responsible_unit_id INT NULL,
    third_party_dependencies NVARCHAR(MAX) NULL,
    outsourcing_partners NVARCHAR(MAX) NULL,
    customer_user NVARCHAR(MAX) NULL,
    process_operation_team_id INT NULL, 
    process_recovery_team_id INT NULL, 
    recovery_manager_id INT NULL,
    recovery_manager_contact_details NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_system_registers_organization_units_1 FOREIGN KEY (business_owner_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_organization_units_2 FOREIGN KEY (it_responsible_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_teams_1 FOREIGN KEY (process_operation_team_id) REFERENCES job_teams (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_teams_2 FOREIGN KEY (process_recovery_team_id) REFERENCES job_teams (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_users_1 FOREIGN KEY (recovery_manager_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_system_registers_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION  
    
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE system_registers;
-- +goose StatementEnd