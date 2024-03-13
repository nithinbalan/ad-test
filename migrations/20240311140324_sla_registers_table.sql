-- +goose Up
-- +goose StatementBegin
CREATE TABLE sla_registers (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_id VARCHAR(45) NOT NULL,
    sla_category_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    responsible_user_id INT NOT NULL,
    lifecycle_status_id INT NOT NULL,
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description VARCHAR(MAX) NULL,
    client VARCHAR(100) NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_sla_registers_sla_categories1 FOREIGN KEY (sla_category_id) REFERENCES sla_categories(id),
    CONSTRAINT FK_sla_registers_severity_levels1 FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id),
    CONSTRAINT FK_sla_registers_responsible_units1 FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id),
    CONSTRAINT FK_sla_registers_responsible_users1 FOREIGN KEY (responsible_user_id) REFERENCES users(id),
    CONSTRAINT FK_sla_registers_lifecycle_statuses1 FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses(id),
    CONSTRAINT FK_sla_registers_created_by1 FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT FK_sla_registers_updated_by1 FOREIGN KEY (updated_by) REFERENCES users(id)

);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.sla_registers;
-- +goose StatementEnd
