-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_register (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    sla_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    lifecycle_status_id INT NOT NULL,
    compliance_status_id INT NULL,
    title_en NVARCHAR(250) NOT NULL,
    title_ar NVARCHAR(250) NOT NULL,
    description NVARCHAR(MAX) NULL,
    client VARCHAR(250) NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_sla_register_sla_categories FOREIGN KEY (sla_category_id) REFERENCES sla_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_severity_levels FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_responsible_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_lifecycle_statuses1 FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_user_2 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_sla_register_user_3 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_register;
-- +goose StatementEnd
